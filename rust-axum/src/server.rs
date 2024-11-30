use axum::{
    extract::Query,
    extract::State,
    http::{Method, StatusCode},
    Json,
    response::IntoResponse,
    Router,
    routing::{get, post},
};
use axum_server::tls_rustls::RustlsConfig;
use serde::{
    Deserialize,
    Serialize
};
use std::{
    io,
    net::ToSocketAddrs,
    sync::{Arc, Mutex},
    time::{Duration, Instant}
};
use tower::ServiceBuilder;
use tower_http::cors::{CorsLayer, Any};
use crate::config;
use crate::store;

const RATE_LIMIT_WIN_SIZE_SEC: Duration = Duration::from_secs(1);

#[derive(Deserialize)]
struct AppGetParams {
    #[serde(default)]
    name: String,
}

#[derive(Deserialize)]
struct AppSetParams {
    #[serde(default)]
    name: String,
}

#[derive(Serialize)]
struct AppDefResponse {
    message: String,
}

#[derive(Serialize)]
struct AppGetResponse {
    message: String,
}

#[derive(Serialize)]
struct AppSetResponse {
    message: String,
}

// Rate limit state
struct RateLimit {
    max_requests: usize,
    window: Duration,
    requests: usize,
    last_reset: Instant,
}

// Fixed window rate limiting
// Time is divided into fixed intervals (windows)
// A counter keeps track of the number of requests within each window
// If the counter exceeds a limit, subsequenct requests within that window are rejected until the window resets
impl RateLimit {
    fn new(max_requests: usize, window: Duration) -> Self {
        Self {
            max_requests,
            window,
            requests: 0,
            last_reset: Instant::now(),
        }
    }

    fn check(&mut self) -> bool {
        let now = Instant::now();
        if now.duration_since(self.last_reset) > self.window {
            self.requests = 0;
            self.last_reset = now;
        }

        if self.requests < self.max_requests {
            self.requests += 1;
            true
        } else {
            false
        }
    }
}

async fn def_handler(
    axum::extract::Extension(rate_limit_state): axum::extract::Extension<Arc<Mutex<RateLimit>>>
) -> impl IntoResponse {
    let mut rate_limit = rate_limit_state.lock().unwrap();
    if !rate_limit.check() {
        let app_resp = AppDefResponse { message: "429 Too Many Requests".to_string(), };
        return (StatusCode::TOO_MANY_REQUESTS, Json(app_resp));
    }
    println!("Default()");
    let app_resp = AppDefResponse { message: "404 Not Found".to_string(), };
    (StatusCode::NOT_FOUND, Json(app_resp))
}

async fn get_handler(
    axum::extract::Extension(rate_limit_state): axum::extract::Extension<Arc<Mutex<RateLimit>>>,
    State(store): State<store::Store>,
    Query(params): Query<AppGetParams>,
) -> impl IntoResponse {
    let mut rate_limit = rate_limit_state.lock().unwrap();
    if !rate_limit.check() {
        let app_resp = AppGetResponse { message: "429 Too Many Requests".to_string(), };
        return (StatusCode::TOO_MANY_REQUESTS, Json(app_resp));
    }
    let name = params.name;
    println!("AppGet(name: \"{}\")", name);
    if name == "" {
        let app_resp = AppGetResponse { message: "400 Bad Request".to_string(), };
        return (StatusCode::BAD_REQUEST, Json(app_resp));
    }
    match store.get_thing(&name) {
        Some(str) => {
            let app_resp = AppGetResponse { message: format!("Hello, {}", str), };
            (StatusCode::OK, Json(app_resp))
        }
        None => {
            let app_resp = AppGetResponse { message: "404 Not Found".to_string(), };
            (StatusCode::NOT_FOUND, Json(app_resp))
        }
    }
}

async fn set_handler(
    axum::extract::Extension(rate_limit_state): axum::extract::Extension<Arc<Mutex<RateLimit>>>,
    State(store): State<store::Store>,
    Query(params): Query<AppSetParams>,
) -> impl IntoResponse {
    let mut rate_limit = rate_limit_state.lock().unwrap();
    if !rate_limit.check() {
        let app_resp = AppSetResponse { message: "429 Too Many Requests".to_string(), };
        return (StatusCode::TOO_MANY_REQUESTS, Json(app_resp));
    }
    let name = params.name;
    println!("AppSet(name: \"{}\")", name);
    if name == "" {
        let app_resp = AppSetResponse { message: "400 Bad Request".to_string(), };
        return (StatusCode::BAD_REQUEST, Json(app_resp));
    }
    match store.set_thing(&name) {
        true => {
            let app_resp = AppSetResponse { message: format!("Hello, {}", name), };
            (StatusCode::OK, Json(app_resp))
        }
        false => {
            let app_resp = AppSetResponse { message: "500 Internal Server Error".to_string(), };
            (StatusCode::INTERNAL_SERVER_ERROR, Json(app_resp))
        }
    }
}

pub async fn run(store: store::Store, insecure: bool) -> io::Result<()> {
    let addr = config::get(config::ADDR_KEY);
    let mut addr_iter = addr.to_socket_addrs()?;
    let addr = addr_iter.next().unwrap();
    let cors_origin_str = config::get(config::CORS_ORIGIN_KEY);
    let cors_origin = match cors_origin_str.parse::<axum::http::HeaderValue>() {
        Ok(val) => val,
        Err(e) => {
            let msg = format!("Failed to parse configuration value: {}: {}", cors_origin_str, e);
            return Err(io::Error::other(msg));
        }
    };
    let cors_layer = CorsLayer::new()
        .allow_origin(cors_origin)
        .allow_methods([Method::GET, Method::POST])
        .allow_headers(Any);
    let req_per_sec_str = config::get(config::REQ_PER_SEC_KEY);
    let req_per_sec = match req_per_sec_str.parse::<usize>() {
        Ok(val) => val,
        Err(e) => {
            let msg = format!("Failed to parse configuration value: {}: {}", req_per_sec_str, e);
            return Err(io::Error::other(msg));
        }
    };
    let rate_limit_state = Arc::new(Mutex::new(RateLimit::new(req_per_sec, RATE_LIMIT_WIN_SIZE_SEC)));
    let rate_limit_layer = axum::extract::Extension(rate_limit_state);
    let middleware = ServiceBuilder::new()
        .layer(cors_layer)
        .layer(rate_limit_layer);
    let app = Router::new()
        .route("/v1/get", get(get_handler))
        .route("/v1/set", post(set_handler))
        .fallback(def_handler)
        .with_state(store)
        .layer(middleware);
    println!("Listening on {}", addr);
    if !insecure {
        let server_cert = config::get(config::CERT_KEY);
        let server_privkey = config::get(config::PRIVKEY_KEY);
        let config = RustlsConfig::from_pem_file(server_cert, server_privkey).await?;
        axum_server::bind_rustls(addr, config)
            .serve(app.into_make_service())
            .await?;
    } else {
        axum_server::bind(addr)
            .serve(app.into_make_service())
            .await?;
    }
    Ok(())
}
