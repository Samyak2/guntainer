package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// interface
//   reference -> docker run <image> <command> <args>
//   guntainer -> gun    run         <command> <args>

func main() {
	if len(os.Args) <= 1 {
		log.Println("Gib command")
		return
	}

	if os.Args[1] == "run" {
		run()
	}
}

func run() {
	log.Println("Hmm, running", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// syscall.Sethostname([]byte("hmm"))
	// syscall.Chroot("/something-nonexistent")
	// syscall.Chdir("/")
	// syscall.Mount("proc", "proc", "proc", 0, "")
	// syscall.Unmount("/proc", 0)

	err := cmd.Run()

	if err != nil {
		log.Fatalln("Error in running command", err)
	}
}
