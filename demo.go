package main

import (
	"fmt"
	"os"
	"os/exec"
)

// go run demo.go run commands arguments
// ou go build demo.go && ./demo run commands arguments
// like docker run [options] image [commands and arguments]
func main() {

	switch os.Args[1] {
	case "run":
		// debug message
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

	// debug message
	fmt.Printf("Rodando %v\n", os.Args[2:])

	//parsing the command
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	doStuff(cmd.Run())
}

func doStuff(err error) {
	if err != nil {
		panic(err)
	}
}
