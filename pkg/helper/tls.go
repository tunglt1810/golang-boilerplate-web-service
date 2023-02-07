package helper

import (
	"crypto/tls"
	"crypto/x509"
)

func ConvertToTLSConfig(cert, key, caCert []byte) (*tls.Config, error) {
	tlsCfg := tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if cert != nil && key != nil && caCert != nil {
		clientCert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsCfg.RootCAs = caCertPool
		tlsCfg.Certificates = []tls.Certificate{clientCert}
	} else {
		tlsCfg.ClientAuth = tls.NoClientCert
	}

	return &tlsCfg, nil
}
