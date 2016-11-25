package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// go run demo-pid.go run image commands arguments
// ou go build demo-pid.go && ./demo-pid run image commands arguments
// like docker run [options] image [commands and arguments]
func main() {

	switch os.Args[1] {
	case "run":
		// debug message
		defer exiting()
		run()
	case "child":
		child()
	default:
		panic("Shit just happened...")
	}
}

func exiting() {
	fmt.Printf("Saindo do Container\n")
}

func run() {
	// fork the process
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// creating new namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	doStuff(cmd.Run())
}

func child() {

	// debug message
	fmt.Printf("Rodando %v as PID %d usando imagem %v\n", os.Args[3:], os.Getpid(), os.Args[2])

	//parsing the command
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// parsing the image to be used and pointing to its file system
	rootfs := "/rootfs-" + os.Args[2]
	doStuff(syscall.Chroot(rootfs))
	doStuff(os.Chdir("/"))

	//mounting /proc
	doStuff(syscall.Mount("proc", "proc", "proc", 0, ""))
	doStuff(cmd.Run())
}

func doStuff(err error) {
	if err != nil {
		panic(err)
	}
}
