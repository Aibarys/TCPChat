package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Mu              sync.Mutex
	clients         = make(map[string]client)
	leaving         = make(chan message)
	messages        = make(chan message)
	newNameChan     = make(chan string)
	oldMessages     = []string{}
	NumberOfClients = 0
)

func HandleConn(conn net.Conn) {
	defer conn.Close()
	content, err := os.OpenFile("linuxlogo.txt", os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		log.Println(err)
	}
	logo, err := io.ReadAll(content)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(conn, "Welcome to TCP-Chat!\n%s", logo)
	fmt.Fprintf(conn, "\n[ENTER YOUR NAME]:")

	name := getName(conn)
	GiveHistory(conn)

	Mu.Lock()
	clients[conn.RemoteAddr().String()] = client{name, conn}
	Mu.Unlock()

	messages <- newMessage(name+" has joined our chat...", conn)

	fmt.Fprintf(conn, "["+time.Now().Format("2006-01-02 15:04:05")+"]["+clients[conn.RemoteAddr().String()].Name+"]"+":")

	input := bufio.NewScanner(conn)

	for input.Scan() {

		if input.Text() == "/change" {
			go nameChanger(name, conn)
			name = <-newNameChan
		}

		fmt.Fprintf(conn, "["+time.Now().Format("2006-01-02 15:04:05")+"]["+clients[conn.RemoteAddr().String()].Name+"]"+":")
		Mu.Lock()
		if isValidTxt(input.Text()) && input.Text() != "/change" && isValidString(input.Text()) {
			messages <- newMessage("["+time.Now().Format("2006-01-02 15:04:05")+"]["+clients[conn.RemoteAddr().String()].Name+"]"+":"+input.Text(), conn)
		}
		Mu.Unlock()

	}

	messages <- newMessage(clients[conn.RemoteAddr().String()].Name+" has left our chat...", conn)

	Mu.Lock()
	delete(clients, conn.RemoteAddr().String())
	Mu.Unlock()

	Mu.Lock()
	NumberOfClients--
	Mu.Unlock()
}

func nameChanger(oldName string, conn net.Conn) {
	fmt.Fprintf(conn, "[ENTER YOUR NEW NAME]:")
	name := getName(conn)
	Mu.Lock()
	clients[conn.RemoteAddr().String()] = client{Name: name, Conn: conn}
	Mu.Unlock()
	fmt.Fprintln(conn, "You have changed the name from "+oldName+" to "+name)
	Mu.Lock()
	messages <- newMessage(oldName+" has changed the name to "+name, conn)
	Mu.Unlock()
	newNameChan <- name
}

func newMessage(msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    msg,
		address: addr,
	}
}

func IsVacantName(name string) bool {
	if len(name) == 0 {
		return false
	}
	Mu.Lock()
	for _, i := range clients {
		if strings.ToLower(i.Name) == strings.ToLower(name) {
			Mu.Unlock()
			return false
		}
	}
	Mu.Unlock()
	return true
}

func getName(conn net.Conn) string {
	nameread := bufio.NewReader(conn)
	name, _ := nameread.ReadString('\n')
	name = strings.TrimSpace(name)
	switch {
	case !IsVacantName(name):
		fmt.Fprintf(conn, "the name is taken, enter another name:")
		name = getName(conn)
		return name
	case !isValidString(name):
		fmt.Fprintf(conn, "the name is invalid, enter a valid name:")
		name = getName(conn)
		return name
	case len([]rune(name)) > 10:
		fmt.Fprintf(conn, "the name is too long. Try 10 characters or less, enter a valid name:")
		name = getName(conn)
		return name
	}
	return name
}

func isValidTxt(input string) bool {
	var count int
	for _, r := range input {
		switch r {
		case ' ', '\t', '\n':
			count++
		}
	}
	if count == len(input) {
		return false
	}
	return true
}

func isValidString(s string) bool {
	s = strings.TrimSuffix(s, "\n")
	rxmsg := regexp.MustCompile("^[\u0400-\u04FF\u0020-\u007F]+$")
	if !rxmsg.MatchString(s) {
		return false
	}
	return true
}
