use clap::Parser;
use std::process::exit;

/// A microservice application with a REST API and an SQLite database
#[derive(Parser)]
struct Cli {
    /// Secure
    #[arg(short, long, default_value_t = false)]
    secure: bool,
    /// Config file name
    #[arg(short = 'c', long = "config")]
    config_file_name: String,
}

#[tokio::main]
async fn main() {
    let opts = Cli::parse();
	if opts.secure {
		println!("TLS enabled")
	} else {
		println!("TLS not enabled")
	}
    if let Err(e) = config::open(&opts.config_file_name) {
        eprintln!("Error: {}", e);
        exit(1);
    }
    let store = store::Store::new().unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
    server::run(store, opts.secure).await.unwrap_or_else(|e| {
        eprintln!("Error: {}", e);
        exit(1);
    });
}

pub mod server;
pub mod store;
pub mod config;
