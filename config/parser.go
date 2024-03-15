package config

type parser struct {
	buffer string

	userFlag    bool
	addressFlag bool

	address string
	user    string
}

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

func (p *parser) getBuffer() string {
	o := p.buffer
	p.buffer = ""
	return o
}

func (p *parser) parseChar(c rune) {
	isSC := p.isSpecialChar(c)

	if !isSC {
		p.buffer = p.buffer + string(c)
	}
}

func parseAddress(address string) EncConfig {
	output := EncConfig{}
	parser := parser{}

	for _, c := range address {
		parser.parseChar(c)
	}

	output.SSHUser = parser.user
	output.SSHAddress = parser.address
	output.RepoPath = parser.buffer

	return output
}
