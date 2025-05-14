package main

import (
	"flag"
	"fmt"
	"os"
)

func exit(str string) {
	fmt.Printf("Error: %s\n", str)
	os.Exit(1)
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [--receive] <host> <filename>...\nNote: Args to be supplied when flag not used\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	_recv := flag.Bool("receive", false, "Receive a file")
	flag.Parse()

	if !*_recv {
		if len(flag.Args()) < 2 {
			flag.Usage()
			os.Exit(1)
		}
		host := flag.Args()[0]
		fmt.Println("Sending to:", host)
		for i := range flag.NArg()-1 {
			filename := flag.Args()[i+1]
			send(filename, host)
		}
		fmt.Println("Sent")
	} else if *_recv {
		fmt.Printf("Awaiting files\r")
		receive()
	} 
}
