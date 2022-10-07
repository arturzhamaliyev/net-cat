package userInterface

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Data struct {
	Message string
	Conn    net.Conn
}

type JoinLeave struct {
	Name      string
	JL        bool
	Conn      net.Conn
	WaitGroup *sync.WaitGroup
}

var Wg sync.WaitGroup

func Chat(c net.Conn, messages chan Data, jlCh chan JoinLeave, name string) error {
	select {
	case jlCh <- JoinLeave{name, true, c, &Wg}:
		Wg.Add(1)
	default:
	}
	Wg.Wait()

	for {
		curTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(c, "[%s][%s]:", curTime, name[:len(name)-1])

		if msg, err := bufio.NewReader(c).ReadString('\n'); err == nil && isValidMsg(msg) {
			info := Data{
				Message: fmt.Sprintf("[%s]:%s", name[:len(name)-1], msg),
				Conn:    c,
			}

			if err := saveLog(fmt.Sprintf("[%s]%s", curTime, info.Message)); err != nil {
				return err
			}

			select {
			case messages <- info:
			default:
				fmt.Printf("couldn't send message from %s\n", name)
			}
		} else if err != nil {
			return err
		}
	}
}

func saveLog(msg string) error {
	file, err := os.OpenFile("./internal/logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(msg); err != nil {
		return err
	}

	return nil
}

func isValidMsg(msg string) bool {
	if msg[:len(msg)-1] == "" || (msg[0] >= 0 && msg[0] <= 31) {
		return false
	}
	return true
}
