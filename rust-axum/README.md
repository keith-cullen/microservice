# Microservice

## Develop the Application

1. create a project

        $ cargo new --bin microservice

2. implement the application

## Build and Run the Application

1. install sqlite3

        $ apt install sqlite3 libsqlite3-dev

2. build the application

        $ cargo build

3. run the application without TLS

        $ ./target/debug/microservice -i -c ../config.yaml

4. run the application with TLS

        $ ./target/debug/microservice -c ../config.yaml
