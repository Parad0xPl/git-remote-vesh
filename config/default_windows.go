//go:build windows

package config

func defaultConfig() EncConfig {
	return EncConfig{
		SSHPort:            22,
		SSHMountPath:       "Z:",
		VeraCryptMountPath: "V:",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
