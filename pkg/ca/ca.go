package ca

/*
 * Copyright 2021 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const (
	DEFAULT_TLS_CERTIFICATE_PATH = "/app/secrets/ca/tls.crt" // use as const
	DEFAULT_TLS_KEY_PATH         = "/app/secrets/ca/tls.key" // use as const

	OPSMX_OID_VALUE = 0x6f706d78 // 31-bit max
)

//
// CA holds the state for the certificate authority.
//
type CA struct {
	config *Config
	caCert tls.Certificate
}

//
// Config holds the filenames for a CA, and has mappings for loading from
// YAML or JSON.
//
type Config struct {
	CACertFile string `yaml:"caCertFile,omitempty" json:"caCertFile,omitempty"`
	CAKeyFile  string `yaml:"caKeyFile,omitempty" json:"caKeyFile,omitempty"`
}

func (c *Config) applyDefaults() {
	if len(c.CACertFile) == 0 {
		c.CACertFile = DEFAULT_TLS_CERTIFICATE_PATH
	}
	if len(c.CAKeyFile) == 0 {
		c.CAKeyFile = DEFAULT_TLS_KEY_PATH
	}
}

func (c *CA) loadCertificate() error {
	caCert, err := tls.LoadX509KeyPair(c.config.CACertFile, c.config.CAKeyFile)
	if err != nil {
		return fmt.Errorf("unable to load CA cetificate or key: %v", err)
	}
	c.caCert = caCert
	return nil
}

//
// LoadCAFromFile will load an existing authority.
//
func LoadCAFromFile(c Config) (*CA, error) {
	c.applyDefaults()

	ca := &CA{
		config: &c,
	}

	err := ca.loadCertificate()
	if err != nil {
		return nil, err
	}
	return ca, nil
}

//
// MakeCAFromData does approximately the same thing as MakeCA() except the CA
// contents are loaded from PEM strings.
//
func MakeCAFromData(certPEM []byte, certPrivKeyPEM []byte) (*CA, error) {
	caCert, err := tls.X509KeyPair(certPEM, certPrivKeyPEM)
	if err != nil {
		return nil, err
	}
	ca := &CA{caCert: caCert}
	return ca, nil
}

//
// GetCACertificate returns the public certificate for the CA.
//
func (c *CA) GetCACertificate() []byte {
	return c.caCert.Certificate[0]
}

func toPEM(data []byte, t string) []byte {
	p := &bytes.Buffer{}
	pem.Encode(p, &pem.Block{
		Type:  t,
		Bytes: data,
	})
	return p.Bytes()
}

//
// MakeCertificateAuthority generates a new certificate authority key, and self-signs it.
//
func MakeCertificateAuthority() ([]byte, []byte, error) {
	now := time.Now().UTC()
	rootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"OpsMX API Forwarder CA"},
			Country:      []string{"DF"},
		},
		NotBefore:             now.Add(-10 * time.Second),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
	}
	empty := []byte{}

	priv, err := rsa.GenerateKey(crand.Reader, 4096)
	if err != nil {
		return empty, empty, err
	}

	// Self-sign the CA key.
	certBytes, err := x509.CreateCertificate(crand.Reader, rootTemplate, rootTemplate, &priv.PublicKey, priv)
	if err != nil {
		return empty, empty, err
	}

	certPEM := toPEM(certBytes, "CERTIFICATE")
	certPrivKeyPEM := toPEM(x509.MarshalPKCS1PrivateKey(priv), "RSA PRIVATE KEY")
	return certPEM, certPrivKeyPEM, nil
}

//
// MakeServerCert will generate a new server certificate, signed with the authority,
// with a validity period of 1 year.  The DNS names will be applied.
//
func (c *CA) MakeServerCert(names []string) (*tls.Certificate, error) {
	now := time.Now().UTC()

	caCert, err := x509.ParseCertificate(c.caCert.Certificate[0])
	if err != nil {
		return nil, err
	}

	certPrivKey, err := rsa.GenerateKey(crand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"OpsMX API Forwarder"},
			Country:      []string{"DF"},
		},
		NotBefore:   now.Add(-10 * time.Second),
		NotAfter:    now.AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    names,
	}

	certBytes, err := x509.CreateCertificate(crand.Reader, certTemplate, caCert, &certPrivKey.PublicKey, c.caCert.PrivateKey)
	if err != nil {
		return nil, err
	}

	certPEM := toPEM(certBytes, "CERTIFICATE")
	certPrivKeyPEM := toPEM(x509.MarshalPKCS1PrivateKey(certPrivKey), "RSA PRIVATE KEY")
	serverCert, err := tls.X509KeyPair(certPEM, certPrivKeyPEM)
	if err != nil {
		return nil, err
	}
	return &serverCert, nil
}

type CertificateName struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Agent   string `json:"agent,omitempty"`
	Purpose string `json:"purpose,omitempty"`
}

const (
	CertificatePurposeControl       = "control"
	CertificatePurposeAgent         = "agent"
	CertificatePurposeService       = "service"
	CertificatePurposeRemoteCommand = "remote-command"
)

func GetCertificateNameFromCert(cert *x509.Certificate) (*CertificateName, error) {
	for _, atv := range cert.Subject.Names {
		if atv.Type.Equal([]int{2, 5, 4, OPSMX_OID_VALUE}) {
			var name CertificateName
			value, ok := atv.Value.(string)
			if !ok {
				return nil, fmt.Errorf("cannot extract custom name from cert, unable to csast to string")
			}
			err := json.Unmarshal([]byte(value), &name)
			if err != nil {
				return nil, err
			}
			return &name, nil
		}
	}
	return nil, fmt.Errorf("did not find custom name in cert")
}

//
// GenerateCertificate will make a new certificate, and return a base64 encoded
// string for the certificate, key, and authority certificate.
//
func (c *CA) GenerateCertificate(name CertificateName) (string, string, string, error) {
	now := time.Now().UTC()
	jsonName, err := json.Marshal(name)
	if err != nil {
		return "", "", "", err
	}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  []int{2, 5, 4, OPSMX_OID_VALUE},
					Value: string(jsonName),
				},
			},
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(crand.Reader, 2048)
	if err != nil {
		return "", "", "", err
	}

	// we now have a certificate and private key.  Now, sign the cert with the CA.

	caCert, err := x509.ParseCertificate(c.caCert.Certificate[0])
	if err != nil {
		return "", "", "", err
	}

	certBytes, err := x509.CreateCertificate(crand.Reader, cert, caCert, &certPrivKey.PublicKey, c.caCert.PrivateKey)
	if err != nil {
		return "", "", "", err
	}

	ca64 := c.GetCACert()
	cert64 := bytesTo64("CERTIFICATE", certBytes)
	certPrivKey64 := bytesTo64("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(certPrivKey))

	return ca64, cert64, certPrivKey64, nil
}

// GetCACert returns the authority certificate encoded as base64.
func (c *CA) GetCACert() string {
	return bytesTo64("CERTIFICATE", c.caCert.Certificate[0])
}

func bytesTo64(prefix string, data []byte) string {
	p := toPEM(data, prefix)
	return base64.StdEncoding.EncodeToString(p)
}

//
// MakeCertPool will return a certificate pool with our CA installed.
//
func (c *CA) MakeCertPool() (*x509.CertPool, error) {
	caCertPool := x509.NewCertPool()
	for _, cert := range c.caCert.Certificate {
		x, err := x509.ParseCertificate(cert)
		if err != nil {
			return nil, fmt.Errorf("unable to parse certificate: %v", err)
		}
		caCertPool.AddCert(x)
		caCertPool.AppendCertsFromPEM(cert)
	}
	return caCertPool, nil
}
