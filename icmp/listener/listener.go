package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"os"
)

const (
	MaxMessageSize = 65535
)

func main() {
	// Start icmp listener, this requires running as privilege.
	conn, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	// Throw an error if the program fail to listen
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, MaxMessageSize)

	f, err := os.Create("res.md")
	if err != nil {
		panic(err)
	}

	for {
		// Read the recieved message
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Extract the message from the request
		msg, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			log.Fatal(err)
		}

		if msg.Type != ipv4.ICMPTypeEcho {
			// Ignore non echo requests
			continue
		}

		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			continue
		}
		payload := echo.Data

		fmt.Printf("Recieved %d bytes\n", len(payload))

		f.Write(payload)

	}
}
