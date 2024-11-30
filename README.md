# Microservice

## Develop, Build, and Run the Application

1. see:
    - go-echo/README.md
    - go-nethttp/README.md
    - rust-actix-web/README.md
    - rust-axum/README.md

## Test the Application using cURL

1. send a GET request

        $ CURL_CA_BUNDLE=./certs/root_server_cert.pem curl -s https://localhost:4443/v1/get?name=Bob | jq

2. send a POST request

        $ CURL_CA_BUNDLE=./certs/root_server_cert.pem curl -s -X POST https://localhost:4443/v1/set?name=Bob | jq

## Test the Application using Postman

on a laptop:

1. install Postman

        goto 'https://www.postman.com/downloads/'

2. configure Postman

        click on 'Import' on the main window
        paste the contents of 'openapi.yaml'
        select 'Postman Collection'

3. configure a request

        select 'Collections'
        select 'App API' -> 'v1' -> 'get' -> 'App Get'
        replace '{{baseUrl}}' with 'https://localhost:4443'
        select 'Params'
        set 'name' to 'Bob'
        select 'Scripts'
        paste
            pm.test("Response status code", function () {
                pm.expect(pm.response.code).to.equal(200);
            });
            pm.test("Response body", function () {
                pm.expect(pm.response.text()).to.include("Hello");
            });

4. send the request

        click on 'Send'

## Test the Application using the Robot Framework

on a laptop:

1. install the Robot Framework

        $ pip install robotframework
        $ pip install robotframework-requests

2. run the test suite

        $ REQUESTS_CA_BUNDLE=./certs/root_server_cert.pem robot test.robot

## Test the Application using RESTler

note: currently, this will only work if the microservice is run without TLS

1. install dependencies

        $ apt-get update
        $ apt-get install -y dotnet-sdk-6.0

2. install restler

        $ mkdir ~/restler_bin
        $ git clone https://github.com/microsoft/restler-fuzzer.git
        $ cd restler-fuzzer
        $ python ./build-restler.py --dest_dir ~/restler_bin

3. generate a grammar

        $ ~/restler_bin/restler/Restler compile --api_spec openapi.yaml

4. run the smoke tests

        $ ~/restler_bin/restler/Restler test \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 4443 \
              --no_ssl

5. run the fuzz-lean tests

        $ ~/restler_bin/restler/Restler fuzz-lean \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 4443 \
              --no_ssl

6. run the fuzz tests

        $ ~/restler_bin/restler/Restler fuzz \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 4443 \
              --no_ssl
              --time_budget=1
