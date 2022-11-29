package main

import (
	"bufio"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"os"
)

const (
	MaxPayloadSize = 65507
)

// noVerifySplitBytes is a simple split function that simply returns its content
func noVerifySplitBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}

// buildICMPEchoRequest build an ICMP echo request with the given data as a byte array
func buildICMPEchoRequest(data []byte) *icmp.Message {
	// Build the echo request
	body := &icmp.Echo{
		ID:   1,
		Seq:  2,
		Data: data,
	}

	// Wrap the echo request in an icmp message that can be sent easily.
	return &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}
}

// sendData build and send an icmp echo request to the given connection with the given bytes
func sendData(conn *icmp.PacketConn, data []byte) {
	log.Printf("Send %d bytes\n", len(data))

	msg := buildICMPEchoRequest(data)

	// Transform the message to a byte array
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteTo(msgBytes, conn.LocalAddr())
	if err != nil {
		panic(err)
	}
}

func main() {

	// Get the destination and filename from the command line arguments
	destination := os.Args[1]
	filename := os.Args[2]

	// Create a connection to the given destination using udp4 to listen with no privilege.
	conn, err := icmp.ListenPacket("udp4", destination)
	if err != nil {
		panic(err)
	}

	// Open the file to send to the listener
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// Close the file when the program end
	defer f.Close()

	// fileReader makes it easier to read the file by bytes
	fileReader := bufio.NewReaderSize(f, MaxPayloadSize)

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(noVerifySplitBytes)

	// Start scanning the file and for each read bytes, send it to the destination.
	for scanner.Scan() {
		sendData(conn, scanner.Bytes())
	}

	// If there has been an error reading the file displays it.
	if err = scanner.Err(); err != nil {
		panic(err)
	}

	log.Println("The file has been sent")
}
