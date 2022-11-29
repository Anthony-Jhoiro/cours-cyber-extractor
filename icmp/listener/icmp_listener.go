package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"os"
)

const (
	MaxIcmpMessageSize   = 65535
	ListeningIcmpNetwork = "ip4:1"
)

func main() {

	// Get the output file from the CLI arguments
	outputFilename := os.Args[1]

	// Open the file, create it if it does not exist and erase its content if it already exists.
	f, err := os.Create(outputFilename)
	if err != nil {
		// If the file counnd't be opened, stop the program
		panic(err)
	}

	// Start icmp listener, this requires running as privilege.
	conn, err := icmp.ListenPacket(ListeningIcmpNetwork, "0.0.0.0")

	// Throw an error if the program fail to listen
	if err != nil {
		panic(err)
	}

	log.Println("Start listening for icmp packets. Stop the program using Ctrl+C when the file is completely sent")

	// this buffer will
	buffer := make([]byte, MaxIcmpMessageSize)

	for {
		// Read the received message
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Extract the message from the request
		msg, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			log.Fatal(err)
		}

		// Ignore non echo requests
		if msg.Type != ipv4.ICMPTypeEcho {
			continue
		}

		// Ignore invalid echo request
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			continue
		}
		payload := echo.Data

		fmt.Printf("Received %d bytes\n", len(payload))

		// Write the message payload in the file
		f.Write(payload)

	}
}
