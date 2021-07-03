package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/mholt/archiver/v3"
)

// interface
//   reference -> docker run  <image>   <command> <args>
//   guntainer -> gun    run  <newroot> <command> <args>
//                       [1]  [2]       [3]

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Gib command")
		return
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	case "childMultiple":
		childMultiple()
	case "build":
		build()
	default:
		log.Fatalln("Unknown command", os.Args[1])
	}
}

func setupNewRoot(archiveName string) (string, error) {
	dname, err := ioutil.TempDir("", "guntainer")
	if err != nil {
		return "", err
	}

	err = archiver.Unarchive(archiveName, dname)

	return dname, err
}

func run() {
	if len(os.Args) <= 2 {
		log.Fatalln("Gib root FS")
	}
	if len(os.Args) <= 3 {
		log.Fatalln("Gib command to run")
	}

	Run(os.Args[2], os.Args[3], os.Args[4:])
}

func Run(imagePath string, command string, args []string) {
	log.Printf("Hmm, running %v with args %v", command, args)

	dname, err := RunNoClean(imagePath, command, args, "child")
	if dname != "" {
		defer os.RemoveAll(dname)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

// runs image
// returns the new root directory, which is not cleaned up
// childCommand must be "child" or "childMultiple"
func RunNoClean(imagePath string, command string, args []string, childCommand string) (string, error) {

	// extract root fs and do stuff
	dname, err := setupNewRoot(imagePath)
	if err != nil {
		return dname, fmt.Errorf("Error extracting root FS: %v", err)
	}
	// log.Println("Extracted root FS to", dname)

	cmd := exec.Command(os.Args[0], append([]string{childCommand, dname, command}, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{}

	// namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      os.Geteuid(),
			Size:        1,
		}},
		GidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      os.Getegid(),
			Size:        1,
		}},
	}

	err = cmd.Run()

	if err != nil {
		return dname, fmt.Errorf("Error running child command: %v", err)
	}

	return dname, nil
}

func child() {
	cmd := exec.Command(os.Args[3], os.Args[4:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname([]byte("guntainer"))
	syscall.Chroot(os.Args[2])
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	err := cmd.Run()

	if err != nil {
		log.Fatalln("Error running command:", err)
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

func childMultiple() {
	// TODO: somehow get multiple command-args here and run them individually
	//       also do the chroot stuff before that

	syscall.Sethostname([]byte("guntainer"))
	syscall.Chroot(os.Args[2])
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	for _, encoded := range os.Args[3:] {
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
