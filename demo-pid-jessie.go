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
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
    }

    doStuff(cmd.Run())
}

func child() {

    fmt.Printf("Rodando %v as PID %d\n", os.Args[2:], os.Getpid())

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    doStuff(syscall.Chroot("/rootfs-jessie"))
    doStuff(os.Chdir("/"))
    doStuff(syscall.Mount("proc", "proc", "proc", 0, ""))
    doStuff(cmd.Run())
}

func doStuff(err error) {
    if err != nil {
        panic(err)
    }
}
