package helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// / Loop contain full logic of handling git communication. Most of logic is
// / based on Alec Newman article (https://rovaughn.github.io/2015-2-9.html)
func Loop() error {
	log.Println("---[Main helper loop]---")
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		command_line, err := stdinReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("can't read next line of communication: %s", err)
		}
		log.Printf("Got command '%s'\n", command_line)

		command_line_parts := strings.Split(command_line, " ")

		command := command_line_parts[0]

		switch {
		case command == "\n":
			return nil
		case command == "capabilities\n":

		case command == "list\n":
		case command == "export\n":
		case command == "capabilities\n":
		case command == "capabilities\n":
		default:
			return fmt.Errorf("unkown command '%s' from git", command)
		}
	}
}
