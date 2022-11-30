package main

import (
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

const Port = 53533

func parseQuery(m *dns.Msg, writer *os.File) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			hexString := strings.TrimSuffix(q.Name, ".")
			bytes, err := hex.DecodeString(hexString)
			if err != nil {
				// ignore b64 decode errors
				continue
			}
			writer.Write(bytes)
		}
	}
}

func makeDnsHandler(file *os.File) func(dns.ResponseWriter, *dns.Msg) {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		if r.Opcode == dns.OpcodeQuery {
			parseQuery(m, file)
		}

		w.WriteMsg(m)
	}
}

func main() {
	f, err := os.Create("res2.md")
	if err != nil {
		panic(err)
	}

	// attach request handler func
	dns.HandleFunc(".", makeDnsHandler(f))

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(Port), Net: "udp"}
	log.Printf("Starting at %d\n", Port)
	err = server.ListenAndServe()

	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}
