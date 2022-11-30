# ICMP Extraction

The objective is to build an extractor using the ICMP protocol.
The emitter and listener are built with GO to allow compilation on almost any OS.

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

# Build the listener (the ./ is very important)
go build ./icmp/listener 

# Build the emitter
go build ./icmp/emiter

# To build the emitter for another platform / arch, you can use GOOS and GOARCH environment variable. For example with windows 32bits :
GOOS=windows GOARCH=386 go build ./icmp/emiter 
# A list of available platform and arch is available here : [https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)
```

## Let's run it !
The emitter executable needs to be sent to the target machine. It doesn't matter how.

run the listener with `./listener`. It requires `sudo` privileges. 






## How it works

### The emitter

The emitter installed of the victim machine will send repetitive requests to the host machine using the ICMP protocol.
Each request will have the following format :
![ICMP request structure](icmp_request_structure.png)

The header is a classical ICMP message header but the data is customised.
It is split into 3 sections.

1. The identifier of the file in 4 bytes. The identifier is generated to favor uniqueness
2. The request index. There is no delay between each request so to ensure the order, we add an index to each message
3. The payload. The par t of the file which is sent in this message.

Whe the whole file is sent, a stop message is sent with the same id as the package, 0 as the request index
and `###STOP###` (encoded in binary) as the payload.

### The listener

The listener listen for any icmp request. For non Stop message, it save the file in the `results` directory. It keeps
one file per request.
Whe the stop request is received. All the files are read and concatenate into a single file that will be located in
the `results` directory.

*Theoretically*, any number of file can be received at the same time as long as they can be handled by the listener.



