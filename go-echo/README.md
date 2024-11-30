# Microservice

## Develop the Application

1. create a project

        $ mkdir microservice
        $ cd microservice
        $ go mod init github.com/keith-cullen/microservice

2. install protoc

        $ cd /tmp
        $ wget https://github.com/protocolbuffers/protobuf/releases/download/v28.0/protoc-28.0-linux-x86_64.zip
        $ unzip protoc-28.0-linux-x86_64.zip
        # cp bin/protoc /usr/local/bin
        # cp -r include/google /usr/local/include

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

8. run entgo

        $ cd store
        $ go run -mod=mod entgo.io/ent/cmd/ent new Thing
        edit ent/schema/thing.go
        $ GOWORK=off go generate ./ent

9. implement the application

## Build and Run the Application

1. build the application

        $ go build

2. run the application without TLS

        $ ./microservice -i -c ../config.yaml

3. run the application with TLS

        $ ./microservice -c ../config.yaml
