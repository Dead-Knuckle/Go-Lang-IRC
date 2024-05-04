# Go-Lang-IRC Client and Server

## Go Chat Client

### Client-Introduction
This project is a terminal-based chat client implemented in Go. It uses the Bubble Tea library to create a user-friendly and visually appealing interface for sending and receiving messages. The client connects to a server via TCP and allows real-time communication with other users connected to the same server.

### Table of Contents
- [Installation](#client-installation)
- [Usage](#client-usage)
- [Features](#client-features)
- [Dependencies](#client-dependencies)
- [Configuration](#client-configuration)

### Client-Installation
To install the chat client, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

Once Go is installed, clone this repository or download the source code:

```bash
git clone https://github.com/Dead-Knuckle/Go-Lang-IRC
cd path-to-your-project
go build -o chat_client
```

### Client-Usage
```bash
./chat_client
```

### Client-Features
- Real-time messaging with other users.
- Beautiful terminal interface using Bubble Tea.
- Connection management over TCP.

### Client-Dependencies
This project depends on several external libraries:

- Bubble Tea
- Lipgloss
- Standard Go libraries for networking and I/O.
- Ensure these dependencies are resolved by running:

```bash
go mod tidy
```

### Client-Configuration
Configuration options include the server address which can be modified in the source code:
```go
const serverAddress = "localhost:6667"
```

## Go Chat Server

### Introduction
This Go-based server application supports the corresponding chat client by handling incoming TCP connections, managing user sessions, and facilitating real-time communication between clients. It features nickname management, heartbeat for connection health checks, and broadcast capabilities.

### Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Dependencies](#dependencies)
- [Configuration](#configuration)
- [Contributors](#contributors)
- [License](#license)

### Server-Installation
To get started with the chat server, ensure that Go is installed on your system. If not, you can download it from [here](https://golang.org/dl/).

After setting up Go, follow these steps:

```bash
git clone https://github.com/Dead-Knuckle/Go-Lang-IRC
cd path-to-your-server-project
go build -o chat_server
```

### Server-Usage
To run the server, execute:

```bash
./chat_server
```
This will start the server on the default configured port (TCP 6667). The server will then be ready to accept connections from clients.

### Server-Features
- Real-Time Communication: Supports messaging between multiple clients in real-time.
- Nickname Management: Allows clients to choose unique nicknames for the session.
- Heartbeat Monitoring: Regularly checks the connection health of clients and disconnects inactive ones.
- Broadcast Messages: Can send messages to all connected clients or specific groups.

###  Server-Dependencies
This project uses the following Go packages:

Standard networking packages (net, bufio).
JSON for message serialization (encoding/json).
These are included in the standard Go library and do not require additional installations.

### Server-Configuration
The server listens on port 6667 by default. You can change the port by modifying the following line in the source code:

```go
listener, err := net.Listen("tcp", ":6667")
```


MIT License

Copyright (c) [2024] [William ELijah Johnson]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
