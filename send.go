package main

import (
	"fmt"
	"net"
	"os"
)

func send(filename, host string) {
	data, err := os.ReadFile(filename)
	length := len(data)
	if err != nil {
		exit(err.Error())
	}

	conn, err := net.Dial("tcp", host+":9191")
	if err != nil {
		exit(err.Error())
	}

	_, err = conn.Write(fmt.Appendf(nil, "%s:%d\n", filename, length))
	if err != nil {
		exit(err.Error())
	}
	_, err = conn.Write(fmt.Appendf(nil, "%s\n", data))
	if err != nil {
		exit(err.Error())
	}
}
