use std::process::exit;

#[tokio::main]
async fn main() {
    let store = store::Store::new().unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
    server::run(store).await.unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
}

pub mod server;
pub mod store;
