use axum::{
    extract::Query,
    extract::State,
    routing::{get, post},
    Json, Router,
};
use axum_server::tls_rustls::RustlsConfig;
use serde::{Deserialize, Serialize};
use std::net::ToSocketAddrs;
use crate::config;
use crate::store;

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
    let response = match store.get_thing(&name) {
        Some(str) => AppGetResponse { message: format!("Hello, {}", str), },
        None => AppGetResponse { message: format!("Internal Server Error"), },
    };
    return Json(response);
}

async fn set_handler(State(store): State<store::Store>, Query(params): Query<AppSetParams>) -> Json<AppSetResponse> {
    let name = params.name;
    println!("AppSet(name: \"{}\")", name);
    let response = match store.set_thing(&name) {
        true => AppSetResponse { message: format!("Hello, {}", name), },
        false => AppSetResponse { message: format!("Internal Server Error"), },
    };
    Json(response)
}

pub async fn run(store: store::Store, secure: bool) -> std::io::Result<()> {
    let server_cert = config::get(config::CERT_KEY);
    let server_privkey = config::get(config::PRIVKEY_KEY);
    let config = RustlsConfig::from_pem_file(server_cert, server_privkey)
        .await?;
    let api_v1 = Router::new()
        .route("/v1/get", get(get_handler))
        .route("/v1/set", post(set_handler));
    let app = Router::new()
        .merge(api_v1)
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
