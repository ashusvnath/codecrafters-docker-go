package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var Debug string = "true"

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	if Debug != "true" {
		log.SetOutput(io.Discard)
	}
	log.Println("Logs from your program will appear here!")

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	dirpath, _ := os.MkdirTemp("", "test-run")
	original_path, err := os.Open(command)
	if err != nil {
		log.Printf("Failed to open original file: %v", err)
		os.Exit(1)
	}
	copied_path, err := os.OpenFile(filepath.Join(dirpath, "executable"), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("Failed to open copy file location: %v", err)
		os.Exit(1)
	}
	io.Copy(copied_path, original_path)
	original_path.Close()
	copied_path.Close()
	all_args := []string{"./executable"}
	all_args = append(all_args, args...)

	syscall.Chroot(dirpath)
	os.Chdir("/")
	cmd := exec.Command("executable", all_args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("executing command %v", cmd)
	os.Mkdir("/dev", 0755)
	devNull, _ := os.Create("/dev/null")
	devNull.Close()
	err = cmd.Run()

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit Status: %d", exiterr.ExitCode())
			os.Exit(exiterr.ExitCode())
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	log.Println("Done")
}
