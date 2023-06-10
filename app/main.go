package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

var Debug string = "false"

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	if Debug != "true" {
		log.SetOutput(io.Discard)
	}
	log.Println("Logs from your program will appear here!")

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)

	errP, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error fetching output pipe: %v", err)
		os.Exit(1)
	}

	outP, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error fetching error pipe: %v", err)
		os.Exit(1)
	}

	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting command: %v", err)
		os.Exit(1)
	}
	done := make(chan bool)
	go transfer(outP, errP, done)

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error waiting for command completion: %v", err)
		os.Exit(1)
	}
	<-done
	log.Println("Done")
}

func transfer(outP, errP io.Reader, done chan<- bool) {
	outData, _ := io.ReadAll(outP)
	errData, _ := io.ReadAll(errP)
	_, err := os.Stdout.Write(outData)
	if err != nil {
		log.Printf("Error copying stdout command: %v", err)
		os.Exit(1)
	}

	_, err = os.Stderr.Write(errData)
	if err != nil {
		log.Printf("Error copying stderr command: %v", err)
		os.Exit(1)
	}
	done <- true
}
