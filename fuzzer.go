package main

import (
	"fmt"
	"fuzzer_go/cmd"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

func main() {
	art, err := os.ReadFile("asciiart.txt")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(art))
	color.Magenta("Scanning...\n")
	time.Sleep(2 * time.Second)
	cmd.Execute()
}
