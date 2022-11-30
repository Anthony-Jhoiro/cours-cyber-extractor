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

//readIcmpEchoRequest waits for the next ICMP message. Parse it and return its data.
func readIcmpEchoRequest(conn *icmp.PacketConn) ([]byte, error) {

	buffer := make([]byte, MaxMessageSize)

	// Read the next received message
	n, _, err := conn.ReadFrom(buffer)
	// If we can't read more message exit the program with an error code of 1
	if err != nil {
		log.Fatal(err)
	}

	// Extract the message from the request
	msg, err := icmp.ParseMessage(1, buffer[:n])
	if err != nil {
		return nil, errors.New("fail to parse icmp message")
	}

	// Ignore non echo requests
	if msg.Type != ipv4.ICMPTypeEcho {
		return nil, errors.New("invalid ICMP request")
	}

	// Cast the request into echo type
	echo, ok := msg.Body.(*icmp.Echo)
	if !ok {
		return nil, errors.New("invalid ICMP request")
	}

	return echo.Data, nil
}

//readTram read the data from a ICMP message and extract each part according to the schema located in the README in the ICMP folder
func readTram(data []byte) (uint32, uint32, []byte) {
	return binary.LittleEndian.Uint32(data[0:4]), binary.LittleEndian.Uint32(data[4:8]), data[8:]
}

//readRequests Listen and handles new requests on the connection
// (this is a blocking function)
func readRequests(conn *icmp.PacketConn) {
	// Listen for requests in an infinite loop
	for {
		// Wait the next icmp message. data contains the request data as a byte array
		data, err := readIcmpEchoRequest(conn)
		// Ignore messages in error to keep only valid messages
		if err != nil {
			continue
		}

		// Decode the tram
		id, seq, payload := readTram(data)

		// Check id the received message is a Stop request. If yes compile the file otherwise, write a temporary file
		if reflect.DeepEqual(payload, icmpcommons.StopSequence) {
			handleStopFile(id)
		} else {
			commons.WriteByteFile(id, seq, payload)
		}
	}
}

//handleStopFile build the file from all the temporary packets
func handleStopFile(id uint32) {
	log.Printf("Recieved Stop file\n")

	// Build the files
	fileName, err := commons.BuildFile(id)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}

	log.Printf("Recieved File %s", fileName)
}

func main() {
	// Start icmp listener, this requires running as privilege.
	conn, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	// Throw an error if the program fail to listen
	if err != nil {
		panic(err)
	}

	// Start listening for requests
	readRequests(conn)
}
