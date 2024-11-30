use axum::{
    extract::Query,
    extract::State,
    routing::{get, post},
    Json, Router,
};
use axum_server::tls_rustls::RustlsConfig;
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
use crate::store;

const SERVER_CERT: &str = "../certs/server_cert.pem";
const SERVER_PRIVKEY: &str = "../certs/server_privkey.pem";
const HOST: [u8; 4] = [0, 0, 0, 0];
const PORT: u16 = 4443;

#[derive(Deserialize)]
struct AppGetParams {
    name: String,
}

#[derive(Deserialize)]
struct AppSetParams {
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

async fn get_handler(State(store): State<store::Store>, Query(params): Query<AppGetParams>) -> Json<AppGetResponse> {
    let name = params.name;
    println!("AppGet(name: \"{}\")", name);
    let mut str = "Invalid".to_string();
    if let Some(name) = store.get_thing(&name) {
        str = name;
    }
    let response = AppGetResponse {
        message: str,
    };
    return Json(response);
}

async fn set_handler(State(store): State<store::Store>, Query(params): Query<AppSetParams>) -> Json<AppSetResponse> {
    let name = params.name;
    println!("AppSet(name: \"{}\")", name);
    if !store.set_thing(&name) {
        let response = AppSetResponse {
            message: "Invalid".to_string(),
        };
        return Json(response);
    }
    let response = AppSetResponse {
        message: format!("Hello, {}", name),
    };
    Json(response)
}

pub async fn run(store: store::Store) -> std::io::Result<()> {
    let config = RustlsConfig::from_pem_file(SERVER_CERT, SERVER_PRIVKEY)
        .await?;
    let api_v1 = Router::new()
        .route("/v1/get", get(get_handler))
        .route("/v1/set", post(set_handler));
    let app = Router::new()
        .merge(api_v1)
        .with_state(store);
    let addr = SocketAddr::from((HOST, PORT));
    println!("TLS enabled");
    println!("Listening on {}", addr);
    axum_server::bind_rustls(addr, config)
        .serve(app.into_make_service())
        .await?;
    Ok(())
}
