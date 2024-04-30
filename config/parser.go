package config

import (
	"log"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// parser handles logic of extracting connection properties from URL
type parser struct {
	buffer string

	userFlag    bool
	addressFlag bool

	address string
	user    string
}

// isSpecialChar check if given character is a special character inside URL
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

// getBuffer return current state of buffer and reset it to empty string.
func (p *parser) getBuffer() string {
	o := p.buffer
	p.buffer = ""
	return o
}

// parseChar add non special character to the buffer
func (p *parser) parseChar(c rune) {
	isSC := p.isSpecialChar(c)

	if !isSC {
		p.buffer = p.buffer + string(c)
	}
}

// parseAddress return VeshConfig extracted from the URL of repo
func parseAddress(address string) VeshConfig {
	address = strings.TrimPrefix(address, "vesh://")

	if utils.IsDebug() {
		log.Println("---Parsing address---")
		log.Println("Address:", address)
	}
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
