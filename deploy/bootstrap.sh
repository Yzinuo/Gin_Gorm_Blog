#! /bin/bash
export WEB_BUILD_CONTEXT="../.."

./clean_docker.sh

cd start
docker-compose up -d --build