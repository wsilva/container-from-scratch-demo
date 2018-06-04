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
	default:
		panic("¯\\_(ツ)_/¯")
	}
}

func run() {
	fmt.Println("--Entrando no conteiner / Get into container--")
	fmt.Printf("--Rodando comando %v / Running command %v --\n", os.Args[2:], os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	doStuff(cmd.Run())
}

func exiting() {
	fmt.Println("--Saindo do conteiner / Exiting container--")
}

func doStuff(err error) {
	if err != nil {
		panic(err)
	}
}
