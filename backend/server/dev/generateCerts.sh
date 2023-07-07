#!/bin/bash

mkdir -p certs

# Generate server certs (provided to Postgres instance)
openssl req -new -x509 -days 365 -nodes -text -out certs/server.crt -keyout certs/server.key -subj "/CN=db"
chmod go-rwx certs/server.key
chmod go-rwx certs/server.crt