package main

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"
)

// go run demo.go run commands arguments
func main() {

    switch os.Args[1] {
    case "run":
        defer exiting()
        run()
    default:
        panic("Shit just happened...")
    }

}

func exiting() {
    fmt.Printf("Saindo do Container\n")
}

func run() {
    fmt.Printf("Rodando %v\n", os.Args[2:])

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS,
    }

    doStuff(cmd.Run())
}

func doStuff(err error) {
    if err != nil {
        panic(err)
    }
}
