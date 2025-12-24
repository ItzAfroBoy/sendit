package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/ItzAfroBoy/mbar"
)

func createConnection(host string) (net.Conn, error) {
	return net.Dial("tcp", host+":9191")
}

func write(filetype, filename string, file io.Reader,  size int64, conn net.Conn, mb *mbar.MBar) error {
	_, err := conn.Write(fmt.Appendf(nil, "%s:%s:%d\n", filetype, filename, size))
	if err != nil {
		return err
	}

	bar := mb.Add(filename, int(size))

	_, err = io.Copy(io.MultiWriter(conn, bar), file)
	if err != nil {
		return err
	}

	return nil
}

func send(filename string, conn net.Conn, mb *mbar.MBar) error {
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return sendDirectory(filename, conn, mb)
	}

	return sendSingleFile(filename, conn, mb)
}

func sendDirectory(dir string, conn net.Conn, mb *mbar.MBar) error {
	zipFile := new(bytes.Buffer)
	w := zip.NewWriter(zipFile)
	if err := w.AddFS(os.DirFS(dir)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return write("zip", dir+".zip", zipFile, int64(zipFile.Len()), conn, mb)
}

func sendSingleFile(filename string, conn net.Conn, mb *mbar.MBar) error {
	var filetype string
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, _ := file.Stat()
	length := stat.Size()
	filename = filepath.Base(filename)
	if strings.HasSuffix(filename, ".zip") {
		filetype = "zip"
	} else {
		filetype = "txt"
	}

	return write(filetype, filename, file, length, conn, mb)
}
