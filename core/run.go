package core

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func setupNewRoot(archiveName string) (string, error) {
	dname, err := ioutil.TempDir("", "guntainer")
	if err != nil {
		return "", err
	}

	err = archiver.Unarchive(archiveName, dname)

	return dname, err
}

func Run(imagePath string, command string, args []string) {
	// log.Printf("Hmm, running %v with args %v", command, args)

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
		return dname, fmt.Errorf("error extracting root FS: %v", err)
	}
	// log.Println("Extracted root FS to", dname)
	// log.Println(os.Args[0], append([]string{childCommand, dname, command}, args...))

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
		return dname, fmt.Errorf("error running child command: %v", err)
	}

	return dname, nil
}
