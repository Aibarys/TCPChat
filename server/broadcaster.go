package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func Broadcaster() {
	for {
		select {
		case msg := <-messages:
			logWriter(msg)
			oldMessages = append(oldMessages, msg.text)
			Mu.Lock()
			for _, conn := range clients {
				if msg.address == conn.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprint(conn.Conn, "\n\033[1A"+"\033[K"+msg.text+"\n["+time.Now().Format("2006-01-02 15:04:05")+"]["+conn.Name+"]"+":")
			}
			Mu.Unlock()
		}
	}
}

func logWriter(txt message) {
	f, err := os.OpenFile("logs.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.WriteString(txt.text + "\n")
}

func GiveHistory(conn net.Conn) {
	Mu.Lock()
	for _, msg := range oldMessages {
		fmt.Fprintln(conn, msg)
	}
	Mu.Unlock()
}
