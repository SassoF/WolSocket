// client.go
package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {
		conn, err := connect()
		if err != nil {
			log.Println(err.Error(), "Connection failed. Another attempt in 3 seocnds")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected", conn.RemoteAddr())
		handleConnection(conn)
	}
}

func connect() (net.Conn, error) {
	return net.Dial("tcp", ":9090")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "shutdown" {
			handleShutdown(conn)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Connection error:", err)
	}
}

func handleShutdown(conn net.Conn) {
	defer conn.Close()
	cmd := exec.Command("sudo", "shutdown", "-h", "now")
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
