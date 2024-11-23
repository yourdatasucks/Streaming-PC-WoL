package main

import (
	"fmt"
	"net"
	"net/http"
	"github.com/mdlayher/wol"
	"strings"
	"golang.org/x/crypto/ssh"
	"os"
)

// parseMAC converts a MAC address string (e.g., "xx:xx:xx:xx:xx:xx") into net.HardwareAddr
func parseMAC(mac string) (net.HardwareAddr, error) {
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid MAC address format")
	}
	macAddr := make([]byte, 6)
	for i, part := range parts {
		byteValue, err := fmt.Sscanf(part, "%02X", &macAddr[i])
		if err != nil || byteValue != 1 {
			return nil, fmt.Errorf("invalid MAC address byte: %v", err)
		}
	}
	return macAddr, nil
}

// Wake function sends a magic Wake-on-LAN packet
func wake(w http.ResponseWriter, r *http.Request) {
	macString := "xx:xx:xx:xx:xx:xx" // Replace with your streaming PC's MAC address
	broadcastAddr := "255.255.255.255:9"
	macAddr, err := parseMAC(macString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing MAC address: %v", err), http.StatusInternalServerError)
		return
	}
	client, err := wol.NewClient()
	if err != nil {
		http.Error(w, "Error creating WOL client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = client.Wake(broadcastAddr, macAddr)
	if err != nil {
		http.Error(w, "Error sending magic packet: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Waking up the gaming PC...")
	fmt.Println("Wake-on-LAN packet sent to:", macString)
}

// Sleep function sends a sleep command using SSH to streaming pc (linux-based)
func sleep(w http.ResponseWriter, r *http.Request) {
    fmt.Println("We entered the sleep function")
	// The IP address and SSH credentials of the machine we want to wake up
	streamingIP := "x.x.x.x" // Replace with streaming pc's local IPV4 address
	username := "my-awesome-box"    // Replace with your pc username

	// Set the path to your private key
	privateKeyPath := "/home/my-awesome-user/.ssh/id_ed25519" // Replace with the correct path to your private key

	// Read the private key file
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading private key from %s: %v", privateKeyPath, err), http.StatusInternalServerError)
		return
	}

	// Create an SSH signer using the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing private key: %v", err), http.StatusInternalServerError)
		return
	}

	// Setup SSH client config with private key authentication
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer), // Use public key authentication
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Skipping host key verification (not recommended for production)
	}

	// Establish SSH connection to PC
	fmt.Println("Attempting to connect to pc via SSH...")
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", streamingIP), config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to PC via SSH: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()
	fmt.Println("SSH connection established to PC:", streamingIP)

	// Create a new session for the sleep command
	session, err := client.NewSession()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating SSH session: %v", err), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	// Run the sleep command on PC
	fmt.Println("Sending sleep command to PC...")
	err = session.Run("sudo systemctl suspend") // Suspend (sleep)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending sleep command: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Sent sleep command to PC...")
	fmt.Println("Successfully sent sleep command to PC.")
}

// Serve web interface
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		<html>
		<head>
			<script src="https://unpkg.com/htmx.org@1.9.3"></script>
			<style>
				body {
					font-family: Arial, sans-serif;
					text-align: center;
					padding: 20px;
					background-color: #121212;
					color: #e0e0e0;
					margin: 0;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
				}
				h1 {
					font-size: 2.5em;
					margin-bottom: 30px;
				}
				.container {
					display: flex;
					flex-direction: column;
					align-items: center;
					justify-content: center;
					width: 100%;
					height: 100%;
				}
				button {
					font-size: 3em;
					padding: 40px;
					margin: 10px;
					min-width: 75vw;
					height: auto;
					max-width: 500px;
					max-height: 400px;
					background-color: #007bff;
					color: white;
					border: none;
					border-radius: 15px;
					cursor: pointer;
					transition: background-color 0.3s;
				}
				button:hover {
					background-color: #0056b3;
				}
				#feedback {
					margin-top: 20px;
					font-size: 1.5em;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Streaming PC Control</h1>
				<button hx-get="/wake" hx-target="#feedback" hx-swap="innerHTML" hx-trigger="click">Wake PC</button>
				<button hx-post="/sleep" hx-target="#feedback" hx-swap="innerHTML" hx-trigger="click" id="sleepButton">Put PC to Sleep</button>
				<div id="feedback"></div>
			</div>
		</body>
		</html>
		`)
	})

	// Handle wake and sleep routes
	http.HandleFunc("/wake", wake)
	http.HandleFunc("/sleep", sleep)

	// Start the server
	fmt.Println("Starting server on http://0.0.0.0:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

