#!/bin/sh

printf "%s" "$GOOGLE_AUTH_JSON" > google-auth.json
resizer \
  -account google-auth.json \
  -bucket resizer \
  -dsn "root:@tcp(mysql:3306)/resizer?charset=utf8&parseTime=True" \
  -host "resizer.storage.googleapis.com"
