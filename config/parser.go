package config

import (
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
)

// parser handles the logic of extracting connection properties from a URL.
type parser struct {
	buffer string

	userFlag    bool
	addressFlag bool

	address string
	user    string
}

// isSpecialChar checks if the given character is a special character inside the URL.
func (p *parser) isSpecialChar(c rune) bool {
	if c == '@' && !p.userFlag {
		p.userFlag = true
		p.user = p.getBuffer()
		return true
	}
	if c == ':' && !p.addressFlag {
		p.addressFlag = true
		p.address = p.getBuffer()
		return true
	}
	return false
}

// getBuffer returns the current state of the buffer and resets it to an empty string.
func (p *parser) getBuffer() string {
	o := p.buffer
	p.buffer = ""
	return o
}

// parseChar adds a non-special character to the buffer.
func (p *parser) parseChar(c rune) {
	isSC := p.isSpecialChar(c)

	if !isSC {
		p.buffer = p.buffer + string(c)
	}
}

// parseAddress returns a VeshConfig extracted from the repository URL.
func parseAddress(address string) VeshConfig {
	address = strings.TrimPrefix(address, "vesh://")

	debug.Println("---Parsing address---")
	debug.Println("Address:", address)

	output := VeshConfig{}
	parser := parser{}

	for _, c := range address {
		parser.parseChar(c)
	}

	output.SSHUser = parser.user
	output.SSHAddress = parser.address
	output.RepoPath = parser.buffer

	return output
}
