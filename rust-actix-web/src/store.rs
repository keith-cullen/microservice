use rusqlite::{Connection, Result, Error, params};
use std::sync::{Arc, Mutex, MutexGuard};
use crate::config;

#[derive(Clone)]
pub struct Store {
    conn: Arc<Mutex<Connection>>,
}

#[derive(Debug)]
struct Thing {
    id: i32,
    name: String,
}

impl Store {
    pub fn new() -> Result<Store> {
        let db_name = config::get(config::DATABASE_FILE_KEY);
        let c = Connection::open(db_name)?;
        c.execute(
            "CREATE TABLE IF NOT EXISTS things (
                        id INTEGER PRIMARY KEY,
                        name TEXT NOT NULL
                )",
            [],
        )?;
        Ok(Store { conn: Arc::new(Mutex::new(c)) } )
    }

    pub fn get_thing(&self, name: &str) -> Option<String> {
        let conn = self.conn.lock().unwrap();
        let op = match self.query_thing(&conn, name) {
            Ok(op) => op,
            Err(e) => {
                println!("Error: Failed to get thing: {}", e);
                return None;
            }
        };
        let thing = match op {
            Some(thing) => thing,
            None => {
                // No entries found or mutliple entries found
                println!("Error: Failed to get thing");
                return None;
            }
        };
        println!("Got thing: \"{}\"", thing.name);
        Some(thing.name)
    }

    pub fn set_thing(&self, name: &str) -> bool {
        let conn = self.conn.lock().unwrap();
        let op = match self.query_thing(&conn, name) {
            Ok(op) => op,
            Err(e) => {
                println!("Error: Failed to set thing: {}", e);
                return false;
            }
        };
        match op {
            None => {
                return self.create_thing(&conn, name);
            }
            Some(thing) => {
                return self.update_thing(&conn, name, thing.id);
            }
        };
    }

    fn query_thing(&self, conn: &MutexGuard<Connection>, name: &str) -> Result<Option<Thing>> {
        let mut stmt = conn.prepare("SELECT id, name FROM things")?;
        let rows = stmt.query_map([], |row| {
            Ok(Thing {
                id: row.get(0)?,
                name: row.get(1)?,
            })
        })?;
        let mut count = 0;
        let mut res = Thing { id: 0, name: "".to_string() };
        for row in rows {
            let thing = row?;
            if name == thing.name {
                res = thing;
                count = count + 1;
            }
        }
        if count < 1 {
            return Ok(None);

        } else if count > 1 {
            return Err(Error::ModuleError("Thing not singular".to_string()));
        }
        Ok(Some(res))
    }

    pub fn create_thing(&self, conn: &MutexGuard<Connection>, name: &str) -> bool {
        if let Err(e) = conn.execute(
            "INSERT INTO things (name) VALUES (?1)",
            params![name],
        ) {
            println!("Error: Failed to create thing: {}", e);
            return false;
        }
        println!("Created thing: \"{}\"", name);
        true
    }

    pub fn update_thing(&self, conn: &MutexGuard<Connection>, name: &str, id: i32) -> bool {
        if let Err(e) = conn.execute(
            "UPDATE things SET name = ?1 WHERE ID = ?2",
            params![name, id],
        ) {
            println!("Error: Failed to update thing: {}", e);
            return false;
        }
        println!("Updated thing: \"{}\"", name);
        true
    }
}
