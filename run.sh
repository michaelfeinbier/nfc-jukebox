#!/bin/sh

. ./.env

concurrently "cd app && yarn dev" "cd server && go run ."
