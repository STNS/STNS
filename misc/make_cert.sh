#!/bin/bash
mkdir demoCA
touch demoCA/index.txt
echo 00 > demoCA/serial
# ca
openssl genrsa -out ./ca.key 2048
openssl req -new -key ca.key -out ca.csr -subj '/C=JP/ST=Tokyo/L=Shibuya-ku/O=Oreore CA inc./OU=Oreore Gr./CN=Oreore CA'
openssl x509 -days 3650 -in ./ca.csr -req -signkey ./ca.key -out ca.crt
# server
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj '/C=JP/ST=Tokyo/L=Tokyo/O=Oreore CA inc./OU=example Gr./CN=server'
yes | openssl ca -config <(cat /usr/local/etc/openssl/openssl.cnf <( printf "\n[usr_cert]\nsubjectAltName=DNS:server,DNS:server")) -keyfile ./ca.key -outdir ./ -cert ca.crt -in server.csr -out server.crt -days 3650

openssl genrsa -out client.key 1024
openssl req -new -key client.key -out client.csr -subj '/C=JP/ST=Tokyo/L=Tokyo/O=Oreore CA inc./OU=client/CN=client'
yes | openssl ca -config <(cat /usr/local/etc/openssl/openssl.cnf <( printf "\n[usr_cert]\nsubjectAltName=DNS:client")) -keyfile ./ca.key -outdir ./ -cert ca.crt -in client.csr -out client.crt -days 3650
