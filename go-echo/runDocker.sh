#!/bin/bash

export REG="localhost:5000"
export IMG="microservice"
export TAG="latest"
export PORT="443"

set -e
docker run -d -p "${PORT}:${PORT}" --mount "type=volume,src=microservice_volume,dst=/data" "${REG}/${IMG}:${TAG}"
