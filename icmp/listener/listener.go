package main

import (
	"encoding/binary"
	"errors"
	"github.com/Anthony-Jhoiro/cyber-extractor/commons"
	icmpcommons "github.com/Anthony-Jhoiro/cyber-extractor/icmp"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"reflect"
)

const (
	MaxMessageSize = 65535
)

func readIcmpEchoRequest(conn *icmp.PacketConn) ([]byte, error) {

	buffer := make([]byte, MaxMessageSize)

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

	if msg.Type != ipv4.ICMPTypeEcho {
		// Ignore non echo requests
		return nil, errors.New("invalid ICMP request")
	}

	// Cast the request into echo type
	echo, ok := msg.Body.(*icmp.Echo)
	if !ok {
		return nil, errors.New("invalid ICMP request")
	}

	return echo.Data, nil
}

func readTram(data []byte) (uint32, uint32, []byte) {
	return binary.LittleEndian.Uint32(data[0:4]), binary.LittleEndian.Uint32(data[4:8]), data[8:]
}

func readRequests(conn *icmp.PacketConn) chan string {
	ch := make(chan string)

	go func() {
		for {
			data, err := readIcmpEchoRequest(conn)
			// Ignore errors
			if err != nil {
				continue
			}

			id, seq, payload := readTram(data)

			if reflect.DeepEqual(payload, icmpcommons.StopSequence) {
				log.Printf("Recieved Stop file\n")

				fileName, err := commons.BuildFile(id)
				if err != nil {
					log.Printf("[ERROR] %v\n", err)
				}

				ch <- fileName
			} else {
				commons.WriteByteFile(id, seq, payload)
			}
		}
	}()

	return ch

}

func main() {
	// Start icmp listener, this requires running as privilege.
	conn, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	// Throw an error if the program fail to listen
	if err != nil {
		panic(err)
	}

	select {
	case file := <-readRequests(conn):
		log.Printf("Recevieved File %s", file)
	}
}
