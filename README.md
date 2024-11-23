# Streaming PC Wake Application

This is a simple Go application designed to wake a streaming PC on the local network using Wake-on-LAN (WOL) and put it to sleep via SSH commands. It was built for a raspberry pi 2 w zero. Additionally, Tailscale is used to access the Raspberry Pi and wake the PC from outside the local network.

## Features
- **Wake Streaming PC Locally**: Sends a Wake-on-LAN magic packet to wake the streaming PC.
- **Put Streaming PC to Sleep**: Uses SSH to issue a suspend command to the streaming PC.
- **Tailscale Integration**: Enables access to the server outside of the local network, ensuring connectivity for waking the PC remotely.
- **Minimal Setup**: Runs a lightweight Go web server on a Raspberry Pi.

## How It Works
1. **Local Network**: Moonlight can send a Wake-on-LAN magic packet to wake the streaming PC while on the local network, but it cannot do so from outside the local network due to the lack of a hardware address to target when using Tailscale.
2. **Outside Network**: Through Tailscale, you can access the Raspberry Pi remotely and trigger the server's wake functionality.
3. **Suspend Command**: When you're done, use the web interface to send a sleep command to the streaming PC via SSH.

## Installation
1. **Clone the Repository**:
   ```bash
   git clone <repository_url>
   cd <repository_name>
   ```

2. **Set Up Your Environment**:
   - Replace placeholders in `main.go` with your:
     - Streaming PC's MAC address.
     - Local IP address of the Streaming PC.
     - SSH username and private key path.

3. **Build the Go Application**:
  [checkout the go docs for more details on building](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
   ```bash
   go build -o streaming_pc_wake main.go
   ```

4. **Run the Server**:
   ```bash
   ./streaming_pc_wake
   ```

5. **Access the Web Interface**:
   - Open a browser and navigate to `http://<raspberry_pi_ip>:5000`.

## Usage
- **Wake Streaming PC**:
  - Click the "Wake PC" button on the web interface to send the magic packet.
- **Put Streaming PC to Sleep**:
  - Click the "Put PC to Sleep" button to send the suspend command via SSH.

## Notes
- **Moonlight Limitation**: Moonlight's magic packet only works when the device is on the same network. This application bridges the gap by allowing WOL functionality via Tailscale.
- **Security**: Ensure SSH is secured with proper key-based authentication, and only trusted devices are on your Tailscale network.

## Troubleshooting
- **WOL Not Working**:
  - Verify the streaming PC's MAC address is correctly set.
  - Ensure the streaming PC is configured for Wake-on-LAN in its BIOS/UEFI.
- **SSH Issues**:
  - Check the private key path and permissions.
  - Test SSH connectivity from the Raspberry Pi to the streaming PC manually.
  - This is where I had the most issues because I was testing on 2 different machines, one worked and the other didn't. It was because I had mistyped or forgotten to change my SSH key path. If you are only viewing the html it will look like it is silently failing if this path is incorrect. Look at the console output for debugging. I may not have caught every fail-state but I caught the ones I had, which will likely be the same that you encounter.
- **Remote Wake Issues**:
  - Confirm Tailscale is properly set up and connected.

