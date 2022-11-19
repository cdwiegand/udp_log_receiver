package main

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// thanks to https://www.golinuxcloud.com/golang-udp-server-client/ for the basic UDP listener code!
var logs = list.New()

func getEnvOr(key string, defaultValue string) (retVal string) {
	retVal, ok := os.LookupEnv(key)
	if !ok {
		retVal = defaultValue
	}
	return
}
func AtoIv2(value string, defaultValue int, minValue int, maxValue int) (retVal int) {
	retVal, err := strconv.Atoi(value)
	if err != nil {
		retVal = defaultValue
	} else if minValue > 0 && retVal < minValue {
		retVal = defaultValue
	} else if maxValue > 0 && retVal > maxValue {
		retVal = defaultValue
	}
	return
}

func main() {
	httpPort := getEnvOr("HTTP_PORT", "8080")
	udpPort := getEnvOr("UDP_PORT", "10000")
	udpBuffer := AtoIv2(getEnvOr("UDP_BUFFER", "65000"), 65000, 1024, 0)
	maxLogLines := AtoIv2(getEnvOr("KEEP_LOGS", "5000"), 5000, 1, 0)
	go runHttpServer(httpPort) // run in background
	go runUdpServer(udpPort, udpBuffer, maxLogLines)
	for {
		time.Sleep(time.Minute)
	}
}

func runUdpServer(udpPort string, udpBuffer int, maxLogLines int) {
	udpServer, err := net.ListenPacket("udp", ":"+udpPort)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, udpBuffer)
		_, _, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		line := string(buf)
		// FIXME: handle multi-line entries
		if logs.Len() >= maxLogLines {
			logs.Remove(logs.Back())
		}
		logs.PushFront(line)
	}
}

func runHttpServer(httpPort string) {
	srv := &http.Server{
		Addr:         ":" + httpPort,
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
	searchFor := []string{}
	if r.URL.Query().Has("q") {
		searchFor = strings.Split(r.URL.Query().Get("q"), ",")
	} else {
		searchFor = nil
	}

	for temp := logs.Front(); temp != nil; temp = temp.Next() {
		if temp.Value != nil {
			tempStr := temp.Value.(string)
			if len(tempStr) > 0 {
				if searchFor != nil {
					for _, searchMe := range searchFor {
						if strings.Contains(tempStr, searchMe) {
							values = append(values, tempStr)
						}
					}
				} else {
					values = append(values, tempStr)
				}
			}
		}
	}
	body := strings.Join(values, "\n") + "\n"

	w.Header().Set("X-Hello", "Darkness, my old friend")
	w.Header().Set("Content-Type", "text/plain")
	bodyBytes := []byte(body)
	w.Header().Set("Content-Length", fmt.Sprint(len(bodyBytes)))
	w.WriteHeader(200)
	w.Write(bodyBytes)
}
