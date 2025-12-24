package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"github.com/ItzAfroBoy/mbar"
)

func exit(str string) {
	fmt.Printf("Error: %s\n", str)
	os.Exit(1)
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [--receive] <host> <filename>...\nNote: Args to be supplied when flag not used\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	recv := flag.Bool("receive", false, "Receive a file")
	flag.Parse()

	mb := mbar.NewMBar()

	if !*recv {
		if len(flag.Args()) < 2 {
			flag.Usage()
			os.Exit(1)
		}
		host := flag.Args()[0]
		fmt.Println("Sending to:", host)
		mb.Start()
		wg := new(errgroup.Group)
		for i := 1; i < flag.NArg(); i++ {
			wg.Go(func() error {
				conn, err := createConnection(host)
				if err != nil {
					return err
				}
				defer conn.Close()
				filename := flag.Args()[i]
				
				err = send(filename, conn, mb)
				if err != nil {
					fmt.Printf("Error sending %s: %v\n", filename, err)
					return err
				}
				return nil
			})
		}
		if err := wg.Wait(); err == nil {
			fmt.Println("All files sent")
		}
	} else {
		fmt.Printf("Awaiting files\r")
		mb.Start()
		receive(mb)
	}
}
