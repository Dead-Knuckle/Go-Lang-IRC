# Go-Lang-IRC Client and Server

## Go Chat Client

### Introduction
This project is a terminal-based chat client implemented in Go. It uses the Bubble Tea library to create a user-friendly and visually appealing interface for sending and receiving messages. The client connects to a server via TCP and allows real-time communication with other users connected to the same server.

### Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Dependencies](#dependencies)
- [Configuration](#configuration)

### Installation
To install the chat client, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

Once Go is installed, clone this repository or download the source code:

```bash
git clone https://your-repository-url.git
cd path-to-your-project
go build -o chat_client
```

### Usage
```bash
./chat_client
```

### Features
- Real-time messaging with other users.
- Beautiful terminal interface using Bubble Tea.
- Connection management over TCP.

### Dependencies
This project depends on several external libraries:

- Bubble Tea
- Lipgloss
- Standard Go libraries for networking and I/O.
- Ensure these dependencies are resolved by running:

```bash
go mod tidy
```

### Configuration
Configuration options include the server address which can be modified in the source code:
```go
const serverAddress = "localhost:6667"
```
