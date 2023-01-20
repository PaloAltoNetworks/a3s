package bootstrap

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"

	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/tg/tglib"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// MakeNATSServer returns an embded nats server. It uses the provide
// conf.NATSConf to configure everything needed. This function WILL change the
// provided NATSConf so it can be used with the generated NATS server.
// Basically, fields in NATSConf related to server connection will be either
// ignored or replaced.
func MakeNATSServer(cfg *conf.NATSConf) (*natsserver.Server, error) {

	// Make a CA
	caCert, caKey, err := tglib.Issue(
		pkix.Name{CommonName: "nats-ca"},
		tglib.OptIssueTypeCA(),
	)

	// Issue a server cert signed by the CA
	serverCert, serverKey, err := tglib.Issue(
		pkix.Name{CommonName: "nats-server"},
		tglib.OptIssueSignerPEMBlock(caCert, caKey, ""),
		tglib.OptIssueIPSANs(net.ParseIP("127.0.0.1")),
		tglib.OptIssueTypeServerAuth(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to issue server cert: %w", err)
	}
	x509ServerKey, err := tglib.PEMToKey(serverKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse server key pem: %w", err)
	}
	x509ServerCert, err := tglib.ParseCertificate(pem.EncodeToMemory(serverCert))
	if err != nil {
		return nil, fmt.Errorf("unable to parse server cert: %w", err)
	}
	tlsServerCert, err := tglib.ToTLSCertificate(x509ServerCert, x509ServerKey)
	if err != nil {
		return nil, fmt.Errorf("unable to convert server cert to tls cert: %w", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem.EncodeToMemory(caCert))

	// Issue a client cert signed by the CA
	clientCert, clientKey, err := tglib.Issue(
		pkix.Name{CommonName: "nats-client"},
		tglib.OptIssueSignerPEMBlock(caCert, caKey, ""),
		tglib.OptIssueTypeClientAuth(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to issue client cert: %w", err)
	}
	x509ClientKey, err := tglib.PEMToKey(clientKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client key pem: %w", err)
	}
	x509ClientCert, err := tglib.ParseCertificate(pem.EncodeToMemory(clientCert))
	if err != nil {
		return nil, fmt.Errorf("unable to parse client cert: %w", err)
	}
	tlsClientCert, err := tglib.ToTLSCertificate(x509ClientCert, x509ClientKey)
	if err != nil {
		return nil, fmt.Errorf("unable to convert client cert to tls cert: %w", err)
	}

	// Instanciate a new nats server with MTLS required.
	nserver, err := natsserver.NewServer(
		&natsserver.Options{
			NoLog:       true,
			Password:    cfg.NATSPassword,
			Username:    cfg.NATSUser,
			Host:        "127.0.0.1",
			AllowNonTLS: false,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{tlsServerCert},
				ClientCAs:    pool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create nats server: %w", err)
	}

	// Adapt conf.NATSConf to allow connecting
	// to the local server.
	cfg.NATSURL = nserver.ClientURL()
	cfg.NATSTLSDisable = false
	cfg.NATSTLSSkip = false
	cfg.NATSCustomTLSConfig = &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{tlsClientCert},
	}

	return nserver, nil
}
