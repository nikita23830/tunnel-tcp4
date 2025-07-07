Tunnel TCP4

[![Release](https://img.shields.io/github/v/release/nikita23830/tunnel-tcp4)](https://github.com/nikita23830/tunnel-tcp4/releases)
[![License](https://img.shields.io/github/license/nikita23830/tunnel-tcp4)](LICENSE)
![Go](https://img.shields.io/badge/made%20with-Go-blue?logo=go)

A lightweight TCP tunneling tool with optional HAProxy PROXY protocol support.  
Easily forward multiple TCP ports to remote servers, and configure everything via a simple JSON file.

## Features

- Forward any local TCP port to a remote server
- Inject HAProxy PROXY-line (optional per tunnel)
- Configuration via a single `config.json` file
- Minimal, fast, and easy to deploy

## Usage

1. **Create `config.json`**

   Example:
   ```json
   [
     {
       "name": "MyTunnel",
       "port": 9000,
       "to": "1.2.3.4:9000",
       "haproxy": true
     }
   ]
2. **Run the binary**

  ```bash
  ./tunnel-tcp4
  ```
  - The program will start all tunnels defined in your config.json.

3. **Check logs**
   You will see messages like:
  ```bash
  [MyTunnel] Tunnel started: :9000 -> 1.2.3.4:9000 (HAProxy: true)
  ```

4. **Download**
  Download prebuilt binaries from the [releases page.](https://github.com/nikita23830/tunnel-tcp4/releases)

