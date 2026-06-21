package main 

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)


func generateTLSConfig() (*tls.Config, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
    return nil, err
}
	template := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
        Organization: []string{"NullBeacon"},
    },
    NotBefore: time.Now(),
    NotAfter:  time.Now().Add(365 * 24 * time.Hour),
    KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
}
	certDER , err := x509.CreateCertificate(rand.Reader,&template, &template, &privateKey.PublicKey, privateKey)
	if err!= nil {
		return nil, err
	}
	
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	
	tlsCert, err := tls.X509KeyPair(certPEM,keyPEM)
	if err!= nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil 
}
