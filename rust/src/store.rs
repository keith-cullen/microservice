use rusqlite::{Connection, Result, params};
use std::sync::{Arc, Mutex, MutexGuard};

const DB_NAME: &str = "store.db";

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
        let c = Connection::open(DB_NAME)?;
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
        let things = match self.query_things(conn) {
            Ok(things) => things,
            Err(e) => {
                println!("Error: Failed to get thing: {}", e);
                return None;
            },
        };
        let mut count = 0;
        let mut res = "".to_string();
        for thing in things {
            println!("id: {}, name: {}", thing.id, thing.name);
            if thing.name == name {
                res = thing.name;
                count = count + 1;
            }
        }
        if count != 1 {
            println!("Error: Failed to get thing: {} entries found", count);
            return None;

        }
        println!("got thing: \"{}\"", res);
        Some(res)
    }

    pub fn set_thing(&self, name: &str) -> bool {
        let conn = self.conn.lock().unwrap();
        return self.create_thing(conn, name)
    }

    fn query_things(&self, conn: MutexGuard<Connection>) -> Result<Vec<Thing>> {
        let mut stmt = conn.prepare("SELECT id, name FROM things")?;
        let rows = stmt.query_map([], |row| {
            Ok(Thing {
                id: row.get(0)?,
                name: row.get(1)?,
            })
        })?;
        let mut things = Vec::new();
        for thing in rows {
            things.push(thing?);
        }
        Ok(things)
    }

    pub fn create_thing(&self, conn: MutexGuard<Connection>, name: &str) -> bool {
        if let Err(e) = conn.execute(
            "INSERT INTO things (name) VALUES (?1)",
            params![name],
        ) {
            println!("Error: Failed to create thing: {}", e);
            return false;
        }
        println!("created thing: \"{}\"", name);
        true
    }

    pub fn update_thing(&self, conn: MutexGuard<Connection>, name: &str) -> bool {
        if let Err(e) = conn.execute(
            "UPDATE things SET name = ?1 WHERE ID = 1",
            params![name],
        ) {
            println!("Error: Failed to update thing: {}", e);
            return false;
        }
        println!("updated thing: \"{}\"", name);
        true
    }
}
