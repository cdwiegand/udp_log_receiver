package main

import (
	"container/list"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// thanks to https://www.golinuxcloud.com/golang-udp-server-client/ for the basic UDP listener code!
var logs = list.New()

func main() {
	go runHttpServer() // run in background
	go runUdpServer()
	for {
		time.Sleep(time.Minute)
	}
}

func runUdpServer() {
	udpServer, err := net.ListenPacket("udp", ":10000")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 65000)
		_, _, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		line := string(buf)
		if logs.Len() >= 5000 {
			logs.Remove(logs.Back())
		}
		logs.PushFront(line)
	}
}

func runHttpServer() {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handleHTTPHandler(handleHTTPRequest),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func handleHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	values := []string{}

	for temp := logs.Front(); temp != nil; temp = temp.Next() {
		values = append(values, temp.Value.(string))
	}
	body := strings.Join(values, "\n")

	w.Header().Set("X-Hello", "Darkness, my old friend")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(body))
}
