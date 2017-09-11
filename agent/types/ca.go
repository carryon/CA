package types

import "crypto/rsa"

type CAConfig struct {
	CA struct {
		Enabled bool `yaml:"enabled"`
		Cert    struct {
			KeyPath string `yaml:"keyPath"`
			CrtPath string `yaml:"crtPath"`
			CaPath  string `yaml:"caPath"`
		}
	}
}

type NodeCert struct {
	PrivateKey      *rsa.PrivateKey
	Certificate     string
	RootCertificate string
}
