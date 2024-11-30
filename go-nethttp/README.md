# Microservice

## Develop the Application

1. create a project

        $ mkdir microservice
        $ cd microservice
        $ go mod init github.com/keith-cullen/microservice

2. run entgo

        $ cd store
        $ go run -mod=mod entgo.io/ent/cmd/ent new Thing
        edit ent/schema/thing.go
        $ GOWORK=off go generate ./ent

3. implement the application

## Build and Run the Application

1. build the application

        $ go build

2. run the application without TLS

        $ ./microservice -i -c ../config.yaml

3. run the application with TLS

        $ ./microservice -c ../config.yaml
