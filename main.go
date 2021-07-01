package main

import (
	"log"
	"io/ioutil"
	"os"
	"os/exec"
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

	log.Println("Hmm, running", os.Args[3:])

	// extract root fs and do stuff
	dname, err := setupNewRoot(os.Args[2])
	if err != nil {
		log.Fatalln("Error extracting root FS:", err)
	}
	// log.Println("Extracted root FS to", dname)
	defer os.RemoveAll(dname)

	cmd := exec.Command(os.Args[0], append([]string{"child", dname}, os.Args[3:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{}

	// namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID: os.Geteuid(),
			Size: 1,
		}},
		GidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID: os.Getegid(),
			Size: 1,
		}},
	}

	err = cmd.Run()

	if err != nil {
		log.Fatalln("Error in running child command:", err)
	}
}

func child() {
	cmd := exec.Command(os.Args[3], os.Args[4:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname([]byte("hmm"))
	syscall.Chroot(os.Args[2])
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	err := cmd.Run()

	if err != nil {
		log.Fatalln("Error in running command:", err)
	}

	syscall.Unmount("/proc", 0)
}
