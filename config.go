package main

type EncConfig struct {
	SSHPort            int
	SSHAddress         string
	SSHUser            string
	SSHIdentityFile    string
	SSHFilePath        string
	SSHMountPath       string
	VeraCryptMountPath string
}

func GetConfig() EncConfig {
	return EncConfig{}
}
