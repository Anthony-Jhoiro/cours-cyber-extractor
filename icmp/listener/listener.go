package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	icmpcommons "github.com/Anthony-Jhoiro/cyber-extractor/icmp"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
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

func buildFile(id uint32) (string, error) {
	tmpFileDir := fmt.Sprintf("results/%d", id)
	fileInfos, err := os.ReadDir(tmpFileDir)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf("results/%d.raw", id))

	defer file.Close()

	sort.Slice(fileInfos, func(i, j int) bool {
		f1, err := strconv.Atoi(fileInfos[i].Name())
		if err != nil {
			return false
		}

		f2, err := strconv.Atoi(fileInfos[j].Name())
		if err != nil {
			return false
		}
		return f1 < f2
	})

	for _, fileInfo := range fileInfos {
		if !fileInfo.Type().IsDir() {
			content, err := os.ReadFile(path.Join(tmpFileDir, fileInfo.Name()))
			if err != nil {
				return "", err
			}
			_, err = file.Write(content)
			if err != nil {
				return "", err
			}
		}
	}

	os.RemoveAll(tmpFileDir)

	return file.Name(), nil
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

				fileName, err := buildFile(id)
				if err != nil {
					log.Printf("[ERROR] %v\n", err)
				}

				ch <- fileName
			} else {
				log.Printf("Recieve [%d] - No %d - %d bytes\n", id, seq, len(payload))

				os.Mkdir(fmt.Sprintf("results/%d", id), 0777)
				file, err := os.Create(fmt.Sprintf("results/%d/%d", id, seq))
				if err != nil {
					panic(err)
				}
				file.Write(payload)

				file.Close()
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
