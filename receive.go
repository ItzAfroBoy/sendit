package main

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strconv"

	"github.com/schollz/progressbar/v3"
)

func receive() {
	ln, err := net.Listen("tcp", ":9191")
	if err != nil {
		exit(err.Error())
	}
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		reader := bufio.NewReader(conn)
		tp := textproto.NewReader(reader)
		buf, _ := tp.ReadLine()
		reg := regexp.MustCompile(`(.+):(\d+)`)
		filename := reg.FindStringSubmatch(string(buf[:]))[1]
		length, _ := strconv.Atoi(reg.FindStringSubmatch(string(buf[:]))[2])
		bar := progressbar.DefaultBytes(int64(length), filename)
		
		file, err := os.Create(filename)
		if err != nil {
			exit(err.Error())
		}

		io.Copy(io.MultiWriter(bar, file), reader)
		bar.Finish()
	}
}
