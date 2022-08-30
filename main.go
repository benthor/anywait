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
	log.SetFlags(log.Llongfile)
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
		log.Printf("started %s", cmd)
		if err != nil {
			cancel()
			log.Println(err)
		}

		go func(cmd *exec.Cmd, cancel context.CancelFunc) {
			err = cmd.Wait()
			log.Printf("%s terminated: %v", cmd, err)
			cancel()
		}(cmd, cancel)
	}

	<-ctx.Done()
}
