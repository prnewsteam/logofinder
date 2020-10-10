#!/bin/bash

cat file.csv | while read line; do
  line=( ${line//,/ } )
  domain="$(echo -e "${line[0]}" | tr -d '[[:space:]]')"
  echo $domain
  curl "http://localhost:8099/logo?domain=$domain&width=120&height=120"
  echo $(ls -1q ../logo/$domain | wc -l)

  if [ $(ls -1q ../logo/$domain | wc -l) -ne "2" ]
  then
    echo "Failed to fetch images from $domain"
  fi
  echo ""
done