package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	macAddress = "XX:XX:XX:XX:XX:XX"

)

type Response struct {
	Message string `json:"message"`
}

type Command struct {
	IP      string `json:"ip"`
	Command string
}

type ConnectionStatus struct {
	net.Conn
	alive bool
}

type Connection struct {
	IP    string `json:"ip"`
	alive bool
}

type Alive struct {
	IP string `json:"ip"`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/index.html")
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/favicon.ico")
}

func PowerOnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := SendWol(macAddress)
	if err != nil {
		json.NewEncoder(w).Encode(Response{err.Error()})
	} else {
		json.NewEncoder(w).Encode(Response{"Package wol sent"})
	}
}

func PowerOffHandler(commandChannel chan<- Command, errorChannel <-chan error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		command := Command{}
		json.NewDecoder(r.Body).Decode(&command)
		command.Command = "shutdown"
		commandChannel <- command

		if err := <-errorChannel; err != nil {
			json.NewEncoder(w).Encode(Response{"No open connection" + err.Error()})
		} else {
			json.NewEncoder(w).Encode(Response{"Shutdown command sent"})
		}

	}

}

func isConnectionAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	buf := make([]byte, 1)
	_, err := conn.Read(buf)

	conn.SetReadDeadline(time.Time{})

	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return true
		}
		return false
	}

	return true
}

func AliveConnectionHandler(informationConChannel chan []ConnectionStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		informationConChannel <- []ConnectionStatus{{nil, false}}
		value := <-informationConChannel
		c := make([]Alive, 0, 5)

		for _, v := range value {
			c = append(c, Alive{v.RemoteAddr().String()})
		}

		if err := json.NewEncoder(w).Encode(c); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}
}

func HandleConnection(connectionChannel <-chan net.Conn, informationConChannel chan []ConnectionStatus,
	errorChannel chan<- error, commandChannel <-chan Command) {
	connection := make([]ConnectionStatus, 0)
	for {
		select {
		case conn := <-connectionChannel:
			if isConnectionAlive(conn) {
				connection = append(connection, ConnectionStatus{conn, true})
			} else {
				for i, v := range connection {
					if v.RemoteAddr().String() == conn.RemoteAddr().String() {
						connection = append(connection[:i], connection[i+1:]...)
					}
				}
			}

		case <-informationConChannel:
			informationConChannel <- connection

		case command := <-commandChannel:
			found := false
			for i, v := range connection {
				if v.RemoteAddr().String() == command.IP {
					found = true
					if isConnectionAlive(connection[i]) {
						_, err := connection[i].Write([]byte(command.Command + "\n"))
						if err != nil {
							errorChannel <- err
						} else {
							errorChannel <- nil
						}
					} else {
						errorChannel <- errors.New("No connection alive")
					}
				}
			}
			if !found {
				errorChannel <- errors.New("No matching connection found for IP: " + command.IP)
			}
		}
	}
}

func HandleInConnection(conn net.Conn, connectionChannel chan<- net.Conn) {
	defer conn.Close()
	connectionChannel <- conn

	for {
		select {
		case <-time.After(3 * time.Second):
			if !isConnectionAlive(conn) {
				connectionChannel <- conn
				log.Println("Connection", conn.RemoteAddr(), "is closed")
				return
			}
		}
	}
}
