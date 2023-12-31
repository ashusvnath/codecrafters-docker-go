package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var Debug string = "true"

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	// if Debug != "true" {
	// 	log.SetOutput(io.Discard)
	// }
	fmt.Println("fmts from your program will appear here!")

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	dirpath, _ := os.MkdirTemp("", "test-run")
	original_path, err := os.Open(command)
	if err != nil {
		fmt.Printf("Failed to open original file: %v", err)
		os.Exit(1)
	}
	copied_path, err := os.OpenFile(filepath.Join(dirpath, "executable"), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Failed to open copy file location: %v", err)
		os.Exit(1)
	}
	io.Copy(copied_path, original_path)
	original_path.Close()
	copied_path.Close()

	syscall.Chroot(dirpath)
	os.Chdir("/")
	wd, _ := os.Getwd()
	fmt.Printf("Current working directory %v\n", wd)
	cmd := exec.Command("./executable", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("executing command %v\n", cmd)
	os.Mkdir("/dev", 0755)
	devNull, _ := os.Create("/dev/null")
	devNull.Close()
	err = cmd.Run()

	if err != nil {
		fmt.Printf("error : %v", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
	fmt.Println("Done")
}
