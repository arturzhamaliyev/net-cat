package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"netcat/internal/userInterface"
)

type User struct {
	Name string
	Conn net.Conn
}

type Channels struct {
	messages chan userInterface.Data
	jlCh     chan userInterface.JoinLeave
}

var (
	count    = 0
	users    = make(map[string]net.Conn)
	channels = ChanConstruct()
)

func ChanConstruct() *Channels {
	return &Channels{make(chan userInterface.Data), make(chan userInterface.JoinLeave)}
}

func Start() error {
	if len(os.Args) > 2 {
		return fmt.Errorf("[USAGE]: ./TCPChat $port")
	}

	HOST, PORT := "localhost:", "8989"
	if len(os.Args) == 2 {
		PORT = os.Args[1]
	}
	fmt.Printf("Listening on the port :%s\n", PORT)

	l, err := net.Listen("tcp", HOST+PORT)
	if err != nil {
		return err
	}
	defer l.Close()

	var m sync.Mutex
	go broadCast(&m)

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go handleConnect(c, &m)
	}
}

func broadCast(m *sync.Mutex) {
	for {
		select {
		case info := <-channels.messages:
			m.Lock()
			for n, c := range users {
				if info.Conn.RemoteAddr() != c.RemoteAddr() {
					realTtime := time.Now().Format("2006-01-02 15:04:05")
					fmt.Fprintf(c, "\n[%s]%s[%s][%s]:", realTtime, info.Message, realTtime, n[:len(n)-1])
				}
			}
			m.Unlock()
		case user := <-channels.jlCh:
			if user.JL {
				fmt.Fprint(user.Conn, printLog())
				user.WaitGroup.Done()
			}

			m.Lock()
			for n, c := range users {
				if user.Conn.RemoteAddr() != c.RemoteAddr() {
					realTtime := time.Now().Format("2006-01-02 15:04:05")
					if user.JL {
						fmt.Fprintf(c, "\n%s has joined our chat...\n[%s][%s]:", user.Name[:len(user.Name)-1], realTtime, n[:len(n)-1])
					} else {
						fmt.Fprintf(c, "\n%s has left our chat...\n[%s][%s]:", user.Name[:len(user.Name)-1], realTtime, n[:len(n)-1])
					}
				}
			}
			m.Unlock()

		}
	}
}

func handleConnect(c net.Conn, m *sync.Mutex) {
	name, err := addUser(c, m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("current users in chat: %d\n", count)

	if err := userInterface.Chat(c, channels.messages, channels.jlCh, name); err != nil {
		if err == io.EOF {
			m.Lock()
			delete(users, name)
			count--
			m.Unlock()
			channels.jlCh <- userInterface.JoinLeave{Name: name, JL: false, Conn: c}
			fmt.Printf("current users in chat: %d\n", count)
			return
		}
		fmt.Println(err)
		return
	}

	c.Close()
}

func printLog() string {
	logs := ""
	file, err := os.Open("./internal/logs.txt")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs += scanner.Text() + "\n"
	}

	return logs
}

func getName(c net.Conn) (string, error) {
	name, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return "", err
	}

	if !isValidName(c, name) {
		return getName(c)
	}

	return name, nil
}

func isValidName(c net.Conn, name string) bool {
	for _, ch := range name[:len(name)-1] {
		if ch >= 0 && ch <= 31 {
			fmt.Fprint(c, "Not valid name\nPlease try again...\n[ENTER YOUR NAME]:")
			return false
		}
	}

	if name[:len(name)-1] == "" {
		fmt.Fprint(c, "Not valid name\nPlease try again...\n[ENTER YOUR NAME]:")
		return false
	}

	if _, ok := users[name]; ok {
		fmt.Fprint(c, "Username has already taken\nPlease try again...\n[ENTER YOUR NAME]:")
		return false
	}

	return true
}

func addUser(c net.Conn, m *sync.Mutex) (string, error) {
	if wlcm, err := userInterface.Welcome(); err == nil {
		fmt.Fprint(c, wlcm[:len(wlcm)-1])
	} else {
		return "", err
	}

	name, err := getName(c)
	if err != nil {
		return "", err
	}

	m.Lock()
	if count >= 10 {
		fmt.Fprint(c, "Lobby is full\n")
		c.Close()
	} else {
		user := User{
			Name: name,
			Conn: c,
		}
		users[user.Name] = user.Conn
		count++
	}
	m.Unlock()

	return name, nil
}
