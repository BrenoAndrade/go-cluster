# README

## Overview
This repository contains a simple implementation of a cluster management system written in Go. The system utilizes TCP connections to communicate between nodes and HTTP endpoints for external interaction.

## Purpose
This repository serves as a learning resource for understanding basic concepts of cluster management, network communication, and HTTP server implementation in Go.

## Functionality
The system consists of two main components:
1. **Cluster**: Manages the collection of nodes and facilitates communication between them.
2. **HTTP Server**: Provides endpoints for adding nodes to the cluster and sending messages to all nodes.

## Usage
1. **Setup**: Ensure you have Go installed on your system.
2. **Clone the Repository**: `git clone https://github.com/brenoandrade/go-cluster.git`
4. **Run the Program**: `PORT=<cluster-port> HTTP=<http-port> go run main.go`
5. **Environment Variables**:
   - `PORT`: Port on which the cluster listens for incoming connections.
   - `HTTP`: Port on which the HTTP server listens for incoming requests.

## Endpoints
- **Add Node**: `/add?addr=<address>` - Add a new node to the cluster.
- **Send Message**: `/send?value=<message>` - Send a message to all nodes in the cluster.

## Example
```bash
# Start the first node
$ PORT=5000 HTTP=6000 go run main.go

# Start the second node
$ PORT=5001 HTTP=6001 go run main.go

# Add a node to the cluster
$ curl localhost:6000/add?addr=localhost:5001

# Send a message to all nodes
$ curl localhost:6000/send?value=hello
```
