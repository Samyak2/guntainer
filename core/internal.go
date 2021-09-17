package core

import (
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func Child(rootPath string, args []string) {
	// log.Println("child", rootPath, args)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname([]byte("guntainer"))
	syscall.Chroot(rootPath)
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	err := cmd.Run()

	if err != nil {
		log.Fatalln("error running command:", err)
	}
}

// each string of cmdWithArgs must be the command and its arguments separated by \u0000
func RunMultipleNoClean(imagePath string, cmdWithArgs []string) (string, error) {
	log.Printf("Hmm, running commands: %v", cmdWithArgs)

	encodedCmdWithArgs := make([]string, len(cmdWithArgs))
	for i, cmd := range cmdWithArgs {
		encoded := base64.StdEncoding.EncodeToString([]byte(cmd))
		encodedCmdWithArgs[i] = encoded
	}

	dname, err := RunNoClean(imagePath, encodedCmdWithArgs[0], encodedCmdWithArgs[1:], "childMultiple")

	return dname, err
}

func ChildMultiple(rootPath string, args []string) {
	syscall.Sethostname([]byte("guntainer"))
	syscall.Chroot(rootPath)
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	for _, encoded := range args {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		cmdWithArgsString := string(decoded)

		divider, err := strconv.Unquote("'\u0000'")
		if err != nil {
			log.Fatalf("Divider could not be unquoted. If you're seeing this error, uhhh. You shoudn't be.")
		}
		cmdWithArgs := strings.Split(cmdWithArgsString, divider)

		if err != nil {
			log.Fatalf("Could not decode command %v: %v", encoded, err)
		}

		log.Printf("Running command %v with args: %v", cmdWithArgs[0], cmdWithArgs[1:])

		cmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()

		if err != nil {
			log.Fatalf("Error running command %v: %v", cmdWithArgs, err)
		}
	}
}
