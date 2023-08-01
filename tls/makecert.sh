#!/bin/bash
# ./makecert.sh some@exmaple.com
rm -rf certs
mkdir certs
openssl req -new -nodes -x509 -sha256 -out certs/server.pem -keyout certs/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
openssl req -new -nodes -x509 -sha256 -out certs/client.pem -keyout certs/client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
