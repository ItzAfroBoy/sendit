package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/ItzAfroBoy/mbar"
)

func receive(mb *mbar.MBar) {
	ln, err := net.Listen("tcp", ":9191")
	if err != nil {
		exit(err.Error())
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go handleConnection(conn, mb)
	}
}

func extractAndWriteFile(base string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	path := filepath.Join(base, f.Name)

	if !strings.HasPrefix(path, filepath.Clean(base)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), f.Mode())
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		if _, err = io.Copy(f, rc); err != nil {
			return err
		}
	}
	return nil
}

func ensureDirectoryExists(filename string) error {
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	return nil
}

func handleConnection(conn net.Conn, mb *mbar.MBar) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	buf, err := tp.ReadLine()
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}

	reg := regexp.MustCompile(`(.+):(.+):(\d+)`)
	matches := reg.FindStringSubmatch(string(buf))
	if len(matches) != 4 {
		fmt.Printf("Invalid header format: %s\n", buf)
		return
	}

	filetype := matches[1]
	filename := matches[2]
	length, _ := strconv.Atoi(matches[3])
	bar := mb.Add(filename, length)

	ensureDirectoryExists(filename)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Create file error: %v\n", err)
		return
	}
	defer file.Close()

	if _, err = io.Copy(io.MultiWriter(file, bar), reader); err != nil {
		fmt.Printf("Copy error: %v\n", err)
		return
	}

	if filetype == "zip" {
		dir := strings.TrimSuffix(filename, ".zip")
		r, err := zip.OpenReader(filename)
		if err != nil {
			fmt.Printf("Open file error: %v\n", err)
			return
		}
		defer r.Close()
		for _, f := range r.File {
			if err := extractAndWriteFile(dir, f); err != nil {
				fmt.Printf("Extract file error: %v\n", err)
				return
			}
		}
	}
}
