//go:build windows

package config

// defaultConfig return default config related to given platform
func defaultConfig() VeshConfig {
	return VeshConfig{
		SSHPort:            22,
		SSHMountPath:       "Z:",
		VeraCryptMountPath: "V:",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
