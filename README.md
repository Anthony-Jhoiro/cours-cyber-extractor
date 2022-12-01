# Extractor tools

This project contains 2 listener and emitters for extraction using ICMP and DNS. they run file with linux listener and linux emitter. Although the go made programms should work on any OS if the Go language can be compile for its platform and arch. See [https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

## Requirements

To build the project, you will need to have go and git installed on your system,
you can then build it using this commands (linux)
```bash
# Clone the project 
git clone https://github.com/Anthony-Jhoiro/cours-cyber-extractor.git

# Move into the repository
cd cours-cyber-extractor

# Download the dependencies
go mod download

# Build the icmp listener (the ./ is very important)
go build -o icmp-listener ./icmp/listener 

# Build the icmp emitter
go build -o icmp-emitter ./icmp/emiter

# Build the dns listener
go build -o dns-emitter ./dns/emiter

# To build an executable for another platform from yours, you can use GOOS and GOARCH environment variable. For example with windows 32bits :
GOOS=windows GOARCH=386 go build ./icmp/emiter 
# A list of available platform and arch is available here : https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
```

## Method details
- [DNS](https://github.com/Anthony-Jhoiro/cours-cyber-extractor/tree/main/dns)
- [ICMP](https://github.com/Anthony-Jhoiro/cours-cyber-extractor/tree/main/icmp)


## FAQ

### Why GO ?

To be able to compile it and run it natively on *any* OS.

### Why make it all so complicated ?

Because I can, and I had fun doing so.
