//go:build linux

package config

// defaultConfig returns the default configuration for the Linux platform.
func defaultConfig() VeshConfig {
	return VeshConfig{
		SSHPort:            22,
		SSHMountPath:       "/tmp/vesh_sshfs",
		VeraCryptMountPath: "/tmp/vesh_veracrypt",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
