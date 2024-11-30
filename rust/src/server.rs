use axum::{
    extract::Query,
    extract::State,
    http::StatusCode,
    Json,
    response::IntoResponse,
    Router,
    routing::{get, post},
};
use axum_server::tls_rustls::RustlsConfig;
use serde::{Deserialize, Serialize};
use std::io;
use std::net::ToSocketAddrs;
use crate::config;
use crate::store;

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
struct AppGetResponse {
    message: String,
}

#[derive(Serialize)]
struct AppSetResponse {
    message: String,
}

async fn def_handler() -> impl IntoResponse {
    println!("Default()");
    let app_resp = Json(AppGetResponse { message: "404 Not Found".to_string(), });
    (StatusCode::NOT_FOUND, app_resp)
}

async fn get_handler(State(store): State<store::Store>, Query(params): Query<AppGetParams>) -> impl IntoResponse {
    let name = params.name;
    println!("AppGet(name: \"{}\")", name);
    if name == "" {
        let app_resp = Json(AppGetResponse { message: "400 Bad Request".to_string(), });
        return (StatusCode::BAD_REQUEST, app_resp);
    }
    match store.get_thing(&name) {
        Some(str) => {
            let app_resp = Json(AppGetResponse { message: format!("Hello, {}", str), });
            (StatusCode::OK, app_resp)
        }
        None => {
            let app_resp = Json(AppGetResponse { message: "404 Not Found".to_string(), });
            (StatusCode::NOT_FOUND, app_resp)
        }
    }
}

async fn set_handler(State(store): State<store::Store>, Query(params): Query<AppSetParams>) -> impl IntoResponse {
    let name = params.name;
    println!("AppSet(name: \"{}\")", name);
    if name == "" {
        let app_resp = Json(AppSetResponse { message: "400 Bad Request".to_string(), });
        return (StatusCode::BAD_REQUEST, app_resp);
    }
    match store.set_thing(&name) {
        true => {
            let app_resp = Json(AppSetResponse { message: format!("Hello, {}", name), });
            (StatusCode::OK, app_resp)
        }
        false => {
            let app_resp = Json(AppSetResponse { message: "500 Internal Server Error".to_string(), });
            (StatusCode::INTERNAL_SERVER_ERROR, app_resp)
        }
    }
}

pub async fn run(store: store::Store, secure: bool) -> io::Result<()> {
    let req_per_sec_str = config::get(config::HTTP_REQ_PER_SEC_KEY);
    let req_per_sec = match req_per_sec_str.parse::<u32>() {
        Ok(rate) => rate,
        Err(e) => {
            let msg = format!("Failed to parse configuration value: {}: {}", req_per_sec_str, e);
            let err = io::Error::new(io::ErrorKind::Other, msg);
            return Err(err);
        }
    };
    let server_cert = config::get(config::CERT_KEY);
    let server_privkey = config::get(config::PRIVKEY_KEY);
    let config = RustlsConfig::from_pem_file(server_cert, server_privkey).await?;
    let app = Router::new()
        .route("/v1/get", get(get_handler))
        .route("/v1/set", post(set_handler))
        .fallback(def_handler)
        .with_state(store);
    if secure {
        let addr = config::get(config::HTTPS_ADDR_KEY);
        let mut addr_iter = addr.to_socket_addrs()?;
        let addr = addr_iter.next().unwrap();
        println!("Listening on {}", addr);
        axum_server::bind_rustls(addr, config)
            .serve(app.into_make_service())
            .await?;
    } else {
        let addr = config::get(config::HTTP_ADDR_KEY);
        let mut addr_iter = addr.to_socket_addrs()?;
        let addr = addr_iter.next().unwrap();
        println!("Listening on {}", addr);
        axum_server::bind(addr)
            .serve(app.into_make_service())
            .await?;
    }
    Ok(())
}
