package main

import (
	"container/list"
	"log"
	"net"
)

// thanks to https://www.golinuxcloud.com/golang-udp-server-client/ for the basic UDP listener code!
var logs = list.New()

func main() {
	udpServer, err := net.ListenPacket("udp", ":10000")
	if err != nil {
		log.Fatal(err)
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
		//newline := fmt.Sprintf("mem: %d: %s", logs.Len(), line)
		//log.Println(newline)
	}
}
