// server.go
package main

import (
	"log"
	"net"
	"net/http"
	"wol/pkg/handlers"
)

func main() {

	commandChannel := make(chan handlers.Command)
	errorChannel := make(chan error)

	connectionChannel := make(chan net.Conn)
	informationConChannel := make(chan []handlers.ConnectionStatus)

	go handlers.HandleConnection(connectionChannel, informationConChannel, errorChannel, commandChannel)

	go func() {
		router := http.NewServeMux()

		router.HandleFunc("GET /{$}", handlers.RootHandler)

		router.HandleFunc("POST /powerOn", handlers.PowerOnHandler)
		router.HandleFunc("POST /shutdownClient", handlers.PowerOffHandler(commandChannel, errorChannel))

		router.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./web/html/"))))
		router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
		router.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("./web/javascript"))))

		router.HandleFunc("GET /favicon.ico", handlers.FaviconHandler)
		router.HandleFunc("GET /aliveConnection", handlers.AliveConnectionHandler(informationConChannel))

		log.Println("Listening on port :8080 for http")
		log.Fatalln(http.ListenAndServe(":8080", router))
	}()

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Listening on port :9090 for socket")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Errore connessione", err.Error())
			continue
		}
		log.Println("Client connesso", conn.RemoteAddr())

		go handlers.HandleInConnection(conn, connectionChannel)
	}

}
