use std::collections::HashMap;
use std::fs;
use std::io::{Result, Error, ErrorKind};
use std::sync::Mutex;
use lazy_static::lazy_static;

pub const DATABASE_FILE_KEY: &str = "DatabaseFile";
pub const CERT_KEY: &str = "Cert";
pub const PRIVKEY_KEY: &str = "Privkey";
pub const HTTP_ADDR_KEY: &str = "HttpAddr";
pub const HTTPS_ADDR_KEY: &str = "HttpsAddr";

lazy_static! {
    static ref MAP: Mutex<HashMap<String, String>> = {
        Mutex::new(HashMap::new())
    };
}

pub fn open(filename: &str) -> Result<()> {
    let content = fs::read_to_string(filename)?;
    let parsed: HashMap<String, String> = match serde_yaml::from_str(&content) {
        Ok(val) => val,
        Err(e) => {
            let msg = format!("Failed to parse configuration file: {}: {}", filename, e);
            let err = Error::new(ErrorKind::Other, msg);
            return Err(err);
        }
    };
    let mut map = MAP.lock().unwrap();
    for (key, value) in &parsed {
        map.insert(key.to_string(), value.to_string());
    }
    Ok(())
}

pub fn get(key: &str) -> String {
    let map = MAP.lock().unwrap();
    match map.get(key) {
        None => {
            println!("Configuration data missing for key: '{}'", key);
            String::from("")
        }
        Some(val) => val.clone()
    }
}
