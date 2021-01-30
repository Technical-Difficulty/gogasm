package main

import (
	"flag"
	"fmt"
	gg "gogasm/src"
)

var header string = `
 ______     ______     ______     ______     ______     __    __    
/\  ___\   /\  __ \   /\  ___\   /\  __ \   /\  ___\   /\ "-./  \   
\ \ \__ \  \ \ \/\ \  \ \ \__ \  \ \  __ \  \ \___  \  \ \ \-./\ \  
 \ \_____\  \ \_____\  \ \_____\  \ \_\ \_\  \/\_____\  \ \_\ \ \_\ 
  \/_____/   \/_____/   \/_____/   \/_/\/_/   \/_____/   \/_/  \/_/ 
`

var seperator string = "--------------------------------------------------------------------"

type config struct {
	wordlist string
	addr     string
	port     string
}

func parseFlags() *config {
	c := &config{}
	flag.StringVar(&c.wordlist, "w", "", "Wordlist file path")
	flag.StringVar(&c.addr, "a", "", "Address eg. localhost")
	flag.StringVar(&c.port, "p", "80", "Port eg. 80")
	flag.Parse()

	return c
}

func printHeader() {
	version := "beta"
	versionLen := len(version) + 2
	seperatorLen := len(seperator)
	seperatorSize := (seperatorLen - versionLen) / 2
	version_seperator := fmt.Sprintf("%s %s %s", seperator[:seperatorSize], version, seperator[:seperatorSize])
	fmt.Printf("%s\n%s\n%s\n\n", version_seperator, header[1:], seperator)
}

func main() {
	printHeader()

	config := parseFlags()

	if config.wordlist == "" {
		fmt.Println("Please provide a wordlist :: $ gogasm -w /path/to/wordlist \n")
		return
	}

	if config.addr == "" {
		fmt.Println("Please provide an address :: $ gogasm -a localhost \n")
		return
	}

	addr := fmt.Sprintf("%s:%s", config.addr, config.port)
	runner := gg.NewRunner("tcp", addr, config.wordlist)
	runner.Start()
}
