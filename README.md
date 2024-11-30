# Microservice

    [https://github.com/google/gnostic/tree/main/cmd/protoc-gen-openapi]
    [https://cloud.google.com/endpoints/docs/grpc/transcoding]
    [https://stackoverflow.com/questions/66168350/import-google-api-annotations-proto-was-not-found-or-had-errors-how-do-i-add]
    [https://github.com/oapi-codegen/oapi-codegen]
    [https://github.com/MarketSquare/robotframework-requests]

## Develop, Build, and Run the Application

1. see:
    - go/README.md
    - rust/README.md

## Test the Applicatioh using cURL

1. send a GET request

        $ CURL_CA_BUNDLE=./certs/root_server_cert.pem curl -s https://localhost:4443/v1/get?name=Bob | jq

2. send a POST request

        $ CURL_CA_BUNDLE=./certs/root_server_cert.pem curl -s -X POST https://localhost:4443/v1/set?name=Bob | jq
