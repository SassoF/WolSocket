package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"wolShutdown/cmd/server/pkg/wol"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/index.html")
}

func powerOnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := wol.SendWol("XX:XX:XX:XX:XX:XX")
	if err != nil {
		json.NewEncoder(w).Encode(Response{err.Error()})
	} else {
		json.NewEncoder(w).Encode(Response{"Package wol sent"})
	}
}

func powerOffHandler(ch chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if conn != nil {
			ch <- "shutdown"
			json.NewEncoder(w).Encode(Response{"Shutdown command sent"})
		} else {
			json.NewEncoder(w).Encode(Response{"No open connection"})
		}
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/css/"+r.PathValue("css"))
}

func cssBootstrapHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/css/bootstrap/"+r.PathValue("css"))
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/javascript/"+r.PathValue("js"))
}

func jsBootstrapHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/javascript/bootstrap/"+r.PathValue("js"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/favicon.ico")
}

func handleConnection(conn net.Conn, ch <-chan string) {
	defer closeConnection()
	msg := <-ch
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func closeConnection() {
	conn.Close()
	conn = nil
}
