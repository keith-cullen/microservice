# Microservice

    [https://github.com/google/gnostic/tree/main/cmd/protoc-gen-openapi]
    [https://cloud.google.com/endpoints/docs/grpc/transcoding]
    [https://stackoverflow.com/questions/66168350/import-google-api-annotations-proto-was-not-found-or-had-errors-how-do-i-add]
    [https://github.com/oapi-codegen/oapi-codegen]
    [https://github.com/MarketSquare/robotframework-requests]

## Develop the Application

1. install protoc

        $ cd /tmp
        $ wget https://github.com/protocolbuffers/protobuf/releases/download/v28.0/protoc-28.0-linux-x86_64.zip
        $ unzip protoc-28.0-linux-x86_64.zip
        # cp bin/protoc /usr/local/bin
        # cp -r include/google /usr/local/include

2. create a project

        $ mkdir microservice
        $ cd microservice
        $ go mod init github.com/keith-cullen/microservice

3. install protoc-gen-openapi

        $ go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

4. download protoc-gen-openapi dependencies

        $ mkdir -p google/api
        goto 'https://github.com/googleapis/googleapis/tree/master/google/api'
        save 'annotations.proto' and 'http.proto' into 'google/api'

        $ mkdir -p google/protobuf
        goto 'https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf'
        save 'wrappers.proto' to 'google/protobuf'

5. run protoc-gen-openapi

        $ protoc app.proto --proto_path=. --openapi_out=.
        see openapi.yaml

6. install oapi-codegen

        $ go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

7. run oapi-codegen

        $ mkdir api
        $ oapi-codegen -config oapi-codegen-config.yaml openapi.yaml

8. implement the application

        edit api/handler.go
        edit main.go

9. run entgo

        $ cd store
        $ go run -mod=mod entgo.io/ent/cmd/ent new Thing
        edit ent/schema/thing.go
        $ GOWORK=off go generate ./ent
        edit store.go

## Build and Run the Application

1. build the application

        $ go build

2. run the application without TLS

        $ ./microservice

3. run the application with TLS

        $ ./microservice -s

4. build a docker image

        $ docker build --tag localhost:5000/microservice:latest .

5. run the docker image

        $ docker run -d -p 8080:8080 localhost:5000/microservice:latest

6. run a kubernetes pod

        $ docker push localhost:5000/microservice:latest
        $ kubectl apply -f pod.yaml

## Test the Applicatioh using cURL

1. send a request

        $ CURL_CA_BUNDLE=./root_server_cert.pem curl https://localhost/v1/get?name=Bob

## Test the Application using the Robot Framework

on a laptop:

1. install the Robot Framework

        $ apt install python3
        $ ln -s /usr/bin/python3 /usr/bin/python
        $ wget https://bootstrap.pypa.io/get-pip.py
        $ chmod +x get-pip.py
        $ ./get-pip.py
        $ pip install --user --upgrade pip
        $ pip install robotframework
        $ pip install robotframework-requests

2. run the test suite

        $ REQUESTS_CA_BUNDLE=./root_server_cert.pem robot test.robot

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
        select 'App API' -> 'v1' -> 'App' -> 'get' -> 'App Get'
        replace '{{baseUrl}}' with 'http://localhost:8080'
        select 'Params'
        set 'name' to 'Bob'
        select 'Scripts'
        paste
            pm.test("Response status code", function () {
                pm.expect(pm.response.code).to.equal(200);
            });
            pm.test("Response body", function () {
                const responseData = pm.response.json();
                pm.expect(responseData).to.be.a('string').and.to.equal("{\"Hello\"}", "Unexpected message body");
            });

4. send the request

        click on 'Send'

## Test the Application using RESTler

note: currently, this will only work if the microservice is built without TLS

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
              --target_port 80 \
              --no_ssl

5. run the fuzz-lean tests

        $ ~/restler_bin/restler/Restler fuzz-lean \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 80 \
              --no_ssl

6. run the fuzz tests

        $ ~/restler_bin/restler/Restler fuzz \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 80 \
              --no_ssl
              --time_budget=1
