package main

import (
	"encoding/hex"
	"github.com/Anthony-Jhoiro/cyber-extractor/commons"
	"log"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

//DnsPort Port for the dns server
const DnsPort = 53533

func handleRecords(m *dns.Msg) {
	// Loop over each record in the dns request
	for _, q := range m.Question {

		// filter on A records
		if q.Qtype != dns.TypeA {
			continue
		}

		// Split the url at the '.' symbol
		parsedName := strings.Split(q.Name, ".")

		// Parse the Url
		payload := parsedName[0]
		sequence, errTs := strconv.ParseUint(parsedName[1], 10, 32)
		id, errId := strconv.ParseUint(parsedName[2], 10, 32)

		// If the payload is "STOP" build the file
		if payload == "STOP" {
			commons.HandleStopFile(uint32(id))
			continue
		}

		// Decode the hexadecimal into binary
		bytes, errB64 := hex.DecodeString(payload)

		// ignore invalid packets
		if errB64 != nil || errTs != nil || errId != nil {
			continue
		}

		// Add the data in a temp file
		commons.WriteByteFile(uint32(id), uint32(sequence), bytes)
	}
}

//dnsHandler handles nex dns requests
func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	// Basic handle of the request to "mock" a real response
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	// handle the records in the request
	if r.Opcode == dns.OpcodeQuery {
		handleRecords(m)
	}

	// Send the response to the emitter
	w.WriteMsg(m)
}

func main() {

	// attach request handler func
	dns.HandleFunc(".", dnsHandler)

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(DnsPort), Net: "udp"}
	log.Printf("Starting at %d\n", DnsPort)

	err := server.ListenAndServe()

	// Shutdown the server when the program exit
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}
