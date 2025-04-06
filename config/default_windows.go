//go:build windows

package config

// defaultConfig returns the default configuration for the Windows platform.
func defaultConfig() VeshConfig {
	return VeshConfig{
		SSHPort:            22,
		SSHMountPath:       "Z:",
		VeraCryptMountPath: "V:",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
