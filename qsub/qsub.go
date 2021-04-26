package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	f = flag.Bool("f", false, "placeholder for qsub foreground flag")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("not enough arguments in call to qsub")
	}
	infile, err := os.Open(args[0])
	defer infile.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "molpro") &&
			!strings.Contains(line, "module") {
			fields := strings.Fields(line)
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			cmd.Run()
		}
	}
	fmt.Printf("%d.maple\n", rand.Intn(1_000_000))
}
