package main

type EncConfig struct {
	SSHPort            int
	SSHAddress         string
	SSHUser            string
	SSHIdentityFile    string
	SSHFilePath        string
	SSHMountPath       string
	VeraCryptMountPath string
	VeraCryptRepoPath  string
}

type SSHFSParams struct {
	SSHPort         int
	SSHAddress      string
	SSHUser         string
	SSHIdentityFile string
	SSHFilePath     string
	SSHMountPath    string
}

func GetConfig() EncConfig {
	return EncConfig{}
}

func (c *EncConfig) ExtractSSHFSParams() SSHFSParams {
	return SSHFSParams{
		SSHPort:         c.SSHPort,
		SSHAddress:      c.SSHAddress,
		SSHUser:         c.SSHUser,
		SSHIdentityFile: c.SSHIdentityFile,
		SSHFilePath:     c.SSHFilePath,
		SSHMountPath:    c.SSHMountPath,
	}
}
