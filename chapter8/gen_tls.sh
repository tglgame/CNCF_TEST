#!/bin/bash

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout grace_http.key -out grace_http.crt -subj "/CN=grace.example.com/O=grace.example.com"

kubectl create secret tls httpsecret --key grace_http.key --cert grace_http.crt -n gracehttp
