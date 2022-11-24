package main

import (
	"container/list"
	"flag"
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
	httpPort := flag.String("http", getEnvOr("HTTP_PORT", "8080"), "HTTP port for API calls")
	udpPort := flag.String("udp", getEnvOr("UDP_PORT", "10000"), "UDP port for receiving logs")
	udpBuffer := flag.String("buf", getEnvOr("UDP_BUFFER", "65000"), "Maximum buffer size for UDP packets")
	maxLogLines := flag.String("lines", getEnvOr("KEEP_LOGS", "5000"), "Maximum number of logs to keep in memory")
	useConsole := flag.Bool("c", truthy(getEnvOr("USE_CONSOLE", "")), "Whether or not to log to console")
	flag.Parse()

	httpPortStr := strconv.Itoa(AtoIv2(*httpPort, 5000, 1, 0))
	udpPortStr := strconv.Itoa(AtoIv2(*udpPort, 5000, 1, 0))
	udpBufferInt := AtoIv2(*udpBuffer, 65000, 1024, 1024*1024*1024*4) // 4 GB should be the max, right?
	maxLogLinesInt := AtoIv2(*maxLogLines, 5000, 1, 0)

	fmt.Println("Using HTTP port", httpPortStr)
	fmt.Println("Using UDP port", udpPortStr, "with a buffer size of", udpBufferInt)
	fmt.Println("Storing", maxLogLinesInt, "log lines at maximum")
	if *useConsole {
		fmt.Println("Printing logs to console")
	}

	go runHttpServer(httpPortStr) // run in background
	go runUdpServer(udpPortStr, udpBufferInt, maxLogLinesInt, *useConsole)
	for {
		time.Sleep(time.Minute)
	}
}

func runUdpServer(udpPort string, udpBuffer int, maxLogLines int, useConsole bool) {
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
		if useConsole {
			fmt.Println(line)
		}
		if logs.Len() >= maxLogLines {
			logs.Remove(logs.Back())
		}
		logs.PushFront(line)
	}
}

func truthy(value string) (ret bool) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	if value == "1" ||
		value == "on" ||
		value == "yes" ||
		value == "y" ||
		value == "true" ||
		value == "t" {
		ret = true
	}
	return
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
