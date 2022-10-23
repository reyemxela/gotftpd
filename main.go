package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/mitchellh/ioprogress"
	"github.com/pin/tftp/v3"
)

func readHandler(filename string, rf io.ReaderFrom) error {
	fmt.Printf("Sending %s\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	finfo, _ := file.Stat()

	reader := &ioprogress.Reader{
		Reader:       file,
		Size:         finfo.Size(),
		DrawInterval: (100 * time.Millisecond),
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, func(progress, total int64) string {
			return fmt.Sprintf("%s %s",
				ioprogress.DrawTextFormatBar(30)(progress, total),
				ioprogress.DrawTextFormatBytes(progress, total),
			)
		}),
	}

	n, err := rf.ReadFrom(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

func writeHandler(filename string, wt io.WriterTo) error {
	fmt.Printf("Receiving %s\n", filename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}

func main() {
	portFlag := flag.Int("p", 69, "port number to use")

	flag.Parse()

	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(5 * time.Second)

	fmt.Printf("Starting server on port %d\n", *portFlag)
	err := s.ListenAndServe(":" + strconv.Itoa(*portFlag))
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}
