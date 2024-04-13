package tool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewReader() *bufio.Reader {
	reader := bufio.NewReader(os.Stdin)
	return reader
}

func ReadStdin(reader *bufio.Reader, describe string) (string, error) {
	fmt.Println(describe)

	readString, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}
	readString = strings.Replace(readString, "\n", "", 1)
	return readString, nil
}


func Input(describe string, defaultString string) (string, error) {
	reader := NewReader()

	readStdin, err := ReadStdin(reader, describe)

	if err != nil {
		LogError.Println(err)
		LogError.Println("read input error")
		return  "", err
	}

	readStdin = strings.Trim(readStdin, "\r")
	if readStdin == "" {
		readStdin = defaultString
	}
	fmt.Println(readStdin)
	return readStdin, nil
}


