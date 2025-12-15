package main

import (
	"fmt"
	"net"
	"os"
)

func addFiles(filename, host string) {
	fmt.Printf(":: %s\n", filename)
	stat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(stat)
	if stat.IsDir() {
		files, err := os.ReadDir(filename)
		if err != nil {
			// exit(err.Error())
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				addFiles(file.Name(), host)
			}
		}
		send(filename, host)
	} else {
		send(filename, host)
	}
}

func send(filename, host string) {
	data, err := os.ReadFile(filename)
	length := len(data)
	if err != nil {
		exit(err.Error())
		panic(err)
	}

	fmt.Println(filename)
	fmt.Println(length)

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
