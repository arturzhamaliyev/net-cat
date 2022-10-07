package userInterface

import (
	"bufio"
	"io"
	"os"
)

func Welcome() (string, error) {
	file, err := os.Open("./internal/userInterface/welcome.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	msg := ""

	reader := bufio.NewReader(file)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			return msg, nil
		} else if err != nil {
			return "", err
		}
		msg += text
	}
}
