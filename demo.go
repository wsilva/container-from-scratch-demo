package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		defer exiting()
		run()
	case "fork":
		fork()
	default:
		panic("¯\\_(ツ)_/¯")
	}
}

func run() {
	fmt.Println("--Entrando no conteiner / Get into container--")
	fmt.Printf("--Imagem usada %v / Image in use %v --\n", os.Args[2], os.Args[2])
	fmt.Printf("--Rodando comando %v / Running command %v --\n", os.Args[3:], os.Args[3:])

	cmd := exec.Command("/proc/self/exe", append([]string{"fork"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	doStuff(cmd.Run())
}

func fork() {

	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	doStuff(syscall.Sethostname([]byte("container")))
	rootfs := "/rootfs-" + os.Args[2]
	doStuff(syscall.Chroot(rootfs))
	doStuff(os.Chdir("/"))
	doStuff(syscall.Mount("proc", "proc", "proc", 0, ""))
	os.Mkdir("/mytemp", 0755)
	doStuff(syscall.Mount("mytemp", "mytemp", "tmpfs", 0, ""))
	doStuff(cmd.Run())
	doStuff(syscall.Unmount("mytemp", 0))
	doStuff(syscall.Unmount("proc", 0))
}

func exiting() {
	fmt.Println("--Saindo do conteiner / Exiting container--")
}

func doStuff(err error) {
	if err != nil {
		panic(err)
	}
}
