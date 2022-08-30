package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())

	var cmd *exec.Cmd

	for _, c := range os.Args[1:] {
		tokens := strings.Split(strings.TrimSpace(c), " ")
		path, err := exec.LookPath(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		cmd = exec.CommandContext(ctx, path, tokens[1:]...)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		//cmd.Stdin = os.Stdin
		err = cmd.Start()
		if err != nil {
			cancel()
			log.Printf("couldn't start %s: %v", cmd, err)
			break
		}
		log.Printf("started %s", cmd)

		go func(cmd *exec.Cmd, cancel context.CancelFunc) {
			err = cmd.Wait()
			log.Printf("%s terminated: %v", cmd, err)
			cancel()
		}(cmd, cancel)
	}

	<-ctx.Done()
}
