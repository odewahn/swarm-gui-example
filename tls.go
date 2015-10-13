package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Reads the specified file into a byte array
func fetchFile(path, fn string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/%s", path, fn)
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error loading %s : %s", fn, err)
	}
	return out, err
}

// GetTLSConfig loads and returns a tls config settings for the swarm.  Lightly adapted version of:
//     https://github.com/ehazlett/interlock/blob/master/interlock/main.go#L14-L32
func GetTLSConfig(certsDir string) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config

	caCert, err := fetchFile(certsDir, os.Getenv("SWARM_CA"))
	cert, err := fetchFile(certsDir, os.Getenv("SWARM_CERT"))
	key, err := fetchFile(certsDir, os.Getenv("SWARM_KEY"))

	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}

	return &tlsConfig, nil
}
