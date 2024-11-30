use actix_cors::Cors;
use actix_governor::{
    Governor,
    GovernorConfigBuilder
};
use actix_web::{
    App,
    get,
    HttpResponse,
    HttpServer,
    post,
    Responder,
    web
};
use serde::Serialize;
use std::{
    fs::File,
    io,
    net::ToSocketAddrs
};
use crate::config;
use crate::store;

const CORS_MAX_AGE_SEC: usize = 3600;

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

async fn def_handler() -> impl Responder {
    let app_resp = AppDefResponse { message: "404 Not Found".to_string(), };
    return HttpResponse::NotFound().json(app_resp);
}

#[get("/v1/get")]
async fn get_handler(
    store: web::Data<store::Store>,
    query: web::Query<std::collections::HashMap<String, String>>
) -> impl Responder {
    let name = match query.get("name") {
        Some(name) => name,
        None => "",
    };
    println!("AppGet(name: \"{}\")", name);
    if name == "" {
        let app_resp = AppGetResponse { message: "400 Bad Request".to_string(), };
        return HttpResponse::BadRequest().json(app_resp);
    }
    match store.get_thing(&name) {
        Some(str) => {
            let app_resp = AppGetResponse { message: format!("Hello, {}", str), };
            return HttpResponse::Ok().json(app_resp);
        }
        None => {
            let app_resp = AppGetResponse { message: "404 Not Found".to_string(), };
            return HttpResponse::NotFound().json(app_resp);
        }
    }
}

#[post("/v1/set")]
async fn set_handler(
    store: web::Data<store::Store>,
    query: web::Query<std::collections::HashMap<String, String>>
) -> impl Responder {
    let name = match query.get("name") {
        Some(name) => name,
        None => "",
    };
    println!("AppSet(name: \"{}\")", name);
    if name == "" {
        let app_resp = AppSetResponse { message: "400 Bad Request".to_string(), };
        return HttpResponse::BadRequest().json(app_resp);
    }
    match store.set_thing(&name) {
        true => {
            let app_resp = AppSetResponse { message: format!("Hello, {}", name), };
            return HttpResponse::Ok().json(app_resp);
        }
        false => {
            let app_resp = AppSetResponse { message: "500 Internal Server Error".to_string(), };
            return HttpResponse::InternalServerError().json(app_resp);
        }
    }
}

pub async fn run(store: store::Store, insecure: bool) -> io::Result<()> {
    let addr = config::get(config::ADDR_KEY);
    let mut addr_iter = addr.to_socket_addrs()?;
    let addr = addr_iter.next().unwrap();
    let cors_origin = config::get(config::CORS_ORIGIN_KEY);
    let req_per_sec_str = config::get(config::REQ_PER_SEC_KEY);
    let req_per_sec = match req_per_sec_str.parse::<u32>() {
        Ok(val) => val,
        Err(e) => {
            let msg = format!("Failed to parse configuration value: {}: {}", req_per_sec_str, e);
            return Err(io::Error::other(msg));
        }
    };
    let governor_conf = match GovernorConfigBuilder::default()
        .requests_per_second(req_per_sec.into())
        .burst_size(req_per_sec)
        .use_headers()
        .finish() {
            Some(val) => val,
            None => {
                let msg = "Failed to configure rate limiter".to_string();
                return Err(io::Error::other(msg));
            }
        };
    let state = web::Data::new(store);
    let app = move || {
        let cors = Cors::default()
            .allowed_origin(&cors_origin)
            .allowed_methods(vec!["GET", "POST"])
            .allowed_headers(vec!["Content-Type", "Authorization"])
            .max_age(CORS_MAX_AGE_SEC);
        let governor = Governor::new(&governor_conf);
        App::new()
            .app_data(state.clone())
            .wrap(governor)
            .wrap(cors)
            .service(set_handler)
            .service(get_handler)
            .default_service(web::route().to(def_handler))
        };
    println!("Listening on {}", addr);
    if !insecure {
        if let Err(_) = rustls::crypto::aws_lc_rs::default_provider()
            .install_default() {
                let msg = "Failed to initialise crypto".to_string();
                return Err(io::Error::other(msg));
            };
        let server_cert = config::get(config::CERT_KEY);
        let server_privkey = config::get(config::PRIVKEY_KEY);
        let mut certs_file = io::BufReader::new(File::open(server_cert)?);
        let mut key_file = io::BufReader::new(File::open(server_privkey)?);
        let tls_certs = match rustls_pemfile::certs(&mut certs_file)
            .collect::<Result<Vec<_>, _>>() {
                Ok(val) => val,
                Err(_) => {
                    let msg = "Failed to parse TLS certificate".to_string();
                    return Err(io::Error::other(msg));
                }
        };
        let tls_key_res = match rustls_pemfile::pkcs8_private_keys(&mut key_file)
            .next() {
                Some(val) => val,
                None => {
                    let msg = "Failed to parse TLS private key".to_string();
                    return Err(io::Error::other(msg));
                }
            };
        let tls_key = match tls_key_res {
            Ok(val) => val,
            Err(e) => {
                let msg = format!("Failed to parse TLS private key: {}", e);
                return Err(io::Error::other(msg));
            }
        };
        let tls_config = match rustls::ServerConfig::builder()
            .with_no_client_auth()
            .with_single_cert(tls_certs, rustls::pki_types::PrivateKeyDer::Pkcs8(tls_key)) {
                Ok(val) => val,
                Err(e) => {
                    let msg = format!("Failed to create TLS configuration: {}", e);
                    return Err(io::Error::other(msg));
                }
            };
        HttpServer::new(app)
            .bind_rustls_0_23(addr, tls_config)?
            .run()
            .await?;
    } else {
        HttpServer::new(app)
            .bind(addr)?
            .run()
            .await?;
    }
    Ok(())
}
