use clap::Parser;
use std::process::exit;

/// A microservice application with a REST API and an SQLite database
#[derive(Parser)]
struct Cli {
    /// Insecure (HTTP) mode
    #[arg(short = 'i', long = "insecure", default_value_t = false)]
    insecure: bool,
    /// Config file name
    #[arg(short = 'c', long = "config")]
    config_file_name: String,
}

#[tokio::main]
async fn main() {
    let opts = Cli::parse();
	if opts.insecure {
		println!("Insecure")
	}
    if let Err(e) = config::open(&opts.config_file_name) {
        eprintln!("Error: {}", e);
        exit(1);
    }
    let store = store::Store::new().unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
    server::run(store, opts.insecure).await.unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
}

pub mod server;
pub mod store;
pub mod config;
