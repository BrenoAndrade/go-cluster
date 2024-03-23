package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	cluster := &Cluster{
		seq:   0,
		sync:  make(chan []byte),
		nodes: []*Node{},
	}

	go setupCluster(cluster)
	go setupHTTP(cluster)

	for {
		select {
		case msg := <-cluster.sync:
			fmt.Println(string(msg))
		}
	}
}

func setupHTTP(cluster *Cluster) {
	port := os.Getenv("HTTP")

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial("tcp", r.URL.Query().Get("addr"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		cluster.AddNode(conn)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")

		payload := fmt.Sprintf("%s %s", port, value)
		cluster.Broadcast(payload)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":"+port, nil)
}

func setupCluster(cluster *Cluster) {
	ln, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %v", err)
		}

		cluster.AddNode(conn)
	}
}

type Node struct {
	id     int
	conn   net.Conn
	active bool
}

func (n *Node) Close() {
	n.conn.Close()
	n.active = false
}

func (n *Node) Write(msg string) {
	_, err := n.conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Error writing node %d: %v", n.id, err)
	}
}

func (n *Node) Listen(sync chan []byte) {
	defer n.Close()

	for {
		buf := make([]byte, 1024)
		_, err := n.conn.Read(buf)
		if err != nil {
			log.Printf("Error reading node %d: %v", n.id, err)
			return
		}

		sync <- append([]byte(fmt.Sprintf("node %d: ", n.id)), buf...)
	}
}

type Cluster struct {
	sync chan []byte

	seq   int
	nodes []*Node
}

func (c *Cluster) AddNode(conn net.Conn) {
	c.seq++
	node := &Node{
		id:     c.seq,
		conn:   conn,
		active: true,
	}
	go node.Listen(c.sync)
	c.nodes = append(c.nodes, node)
}

func (c *Cluster) RemoveNode(id int) {
	if id > len(c.nodes) {
		return
	}

	c.nodes[id].active = false
}

func (c *Cluster) Broadcast(payload string) {
	fmt.Println(payload)

	for _, node := range c.nodes {
		if node.active {
			node.Write(payload)
		}
	}
}

func (c *Cluster) Close() {
	for _, node := range c.nodes {
		node.Close()
	}
}
