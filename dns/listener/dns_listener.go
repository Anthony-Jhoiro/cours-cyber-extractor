package main

import (
	"encoding/hex"
	"github.com/Anthony-Jhoiro/cyber-extractor/commons"
	"log"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

const Port = 53533

func parseQuery(m *dns.Msg, ch chan string) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			parsedName := strings.Split(q.Name, ".")

			payload := parsedName[0]
			sequence, errTs := strconv.ParseUint(parsedName[1], 10, 32)
			id, errId := strconv.ParseUint(parsedName[2], 10, 32)

			// If the payload is "STOP" build the file
			if payload == "STOP" {
				file, err := commons.BuildFile(uint32(id))
				if err == nil {
					ch <- file
				}
				continue
			}

			bytes, errB64 := hex.DecodeString(payload)

			if errB64 != nil || errTs != nil || errId != nil {
				// ignore invalid packets
				continue
			}

			commons.WriteByteFile(uint32(id), uint32(sequence), bytes)
		}
	}
}

func makeDnsHandler(ch chan string) func(dns.ResponseWriter, *dns.Msg) {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		if r.Opcode == dns.OpcodeQuery {
			parseQuery(m, ch)
		}

		w.WriteMsg(m)
	}
}

func main() {

	ch := make(chan string)

	// attach request handler func
	dns.HandleFunc(".", makeDnsHandler(ch))

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(Port), Net: "udp"}
	log.Printf("Starting at %d\n", Port)
	go func() {
		select {
		case file := <-ch:
			log.Println(file)
		}
	}()

	err := server.ListenAndServe()

	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}
