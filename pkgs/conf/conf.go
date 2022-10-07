package conf

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/tg/tglib"
)

// HTTPTimeoutsConf holds http server timeout.
type HTTPTimeoutsConf struct {
	TimeoutRead  time.Duration `mapstructure:"timeout-read"              desc:"Read timeout for the http requests"         default:"120s"`
	TimeoutWrite time.Duration `mapstructure:"timeout-write"             desc:"Write timeout for the http requests"        default:"240s"`
	TimeoutIdle  time.Duration `mapstructure:"timeout-idle"              desc:"Idle timeout for the http requests"         default:"240s"`
}

// LoggingConf is the configuration for log.
type LoggingConf struct {
	LogFormat    string `mapstructure:"log-format" desc:"Log format"                    default:"json"`
	LogLevel     string `mapstructure:"log-level"  desc:"Log level"                     default:"info"`
	LogTracerURL string `mapstructure:"log-tracer" desc:"url of opentracing collector"`
}

// RateLimitingConf holds the configuration for rate limiting.
type RateLimitingConf struct {
	RateLimitingBurst   int  `mapstructure:"rate-limit-burst"   desc:"Burst value"                                 default:"500"`
	RateLimitingEnabled bool `mapstructure:"rate-limit-enabled" desc:"Enable global rate limiting"`
	RateLimitingRPS     int  `mapstructure:"rate-limit-rps"     desc:"Requests per seconds"                        default:"2000"`
}

// ProfilingConf holds the configuration for profiling.
type ProfilingConf struct {
	ProfilingListenAddress string `mapstructure:"profiling-listen"      desc:"Listening address for the profiling server"  default:":6060"`
	ProfilingEnabled       bool   `mapstructure:"profiling-enabled"     desc:"Enable the profiling server"`
}

// HealthConfiguration holds the configuration for health.
type HealthConfiguration struct {
	HealthListenAddress string `mapstructure:"health-listen"            desc:"Listening address for the health server"     default:":80"`
	EnableHealth        bool   `mapstructure:"health-enabled"           desc:"Enable the health check server"`
}

// APIServerConf holds the basic server conf.
type APIServerConf struct {
	PublicAPIURL          string   `mapstructure:"public-api-url"            desc:"Publicly announced API URL"`
	CORSAdditionalOrigins []string `mapstructure:"cors-additional-origins"   desc:"Set additional allowed origin for CORS"`
	CORSDefaultOrigin     string   `mapstructure:"cors-default-origin"       desc:"Set the default allowed origin for CORS"`
	ListenAddress         string   `mapstructure:"listen"                    desc:"Listening address"                                    default:":443"`
	MaxConnections        int      `mapstructure:"max-conns"                 desc:"Max number concurrent TCP connection"`
	MaxProcs              int      `mapstructure:"max-procs"                 desc:"Set the max number thread Go will start"`
	TLSCertificate        string   `mapstructure:"tls-cert"                  desc:"Path to the certificate for https"`
	TLSClientCA           string   `mapstructure:"tls-client-ca"             desc:"Path to the CA to use to verify client certificates"`
	TLSDisable            bool     `mapstructure:"tls-disable"               desc:"Completely disable TLS support"`
	TLSKey                string   `mapstructure:"tls-key"                   desc:"Path to the key for https"`
	TLSKeyPass            string   `mapstructure:"tls-key-pass"              desc:"Password for the key"                                 secret:"true" file:"true"`
}

// TLSConfig returns the configured TLS configuration as *tls.Config.
func (c *APIServerConf) TLSConfig() (*tls.Config, error) {

	if c.TLSDisable {
		return nil, nil
	}

	tlscfg := &tls.Config{}

	if c.TLSClientCA != "" {
		caData, err := os.ReadFile(c.TLSClientCA)
		if err != nil {
			return nil, fmt.Errorf("unable to load ca file: %w", err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		tlscfg.ClientCAs = pool
	}

	if c.TLSCertificate != "" {
		cert, key, err := tglib.ReadCertificatePEM(c.TLSCertificate, c.TLSKey, c.TLSKeyPass)
		if err != nil {
			return nil, fmt.Errorf("unable to load client certificate: %w", err)
		}
		tlscert, err := tglib.ToTLSCertificate(cert, key)
		if err != nil {
			return nil, fmt.Errorf("unable to convert to tls.Certificate: %w", err)
		}
		tlscfg.Certificates = []tls.Certificate{tlscert}
	}

	return tlscfg, nil
}

// MongoConf holds the configuration for mongo db authentication.
type MongoConf struct {
	MongoAttrEncryptKey string `mapstructure:"mongo-encryption-key" desc:"Key to use for attributes encryption"         secret:"true" file:"true"`
	MongoAuthDB         string `mapstructure:"mongo-auth-db"        desc:"Database to use for authenticating"           default:"admin"`
	MongoConsistency    string `mapstructure:"mongo-consistency"    desc:"Set the read consistency"                     default:"nearest" allowed:"strong,monotonic,eventual,nearest,weakest"`
	MongoDBName         string `mapstructure:"mongo-db"             desc:"Database name in MongoDB"                     default:"override-me"`
	MongoPassword       string `mapstructure:"mongo-pass"           desc:"Password to use to connect to Mongo"          secret:"true" file:"true"`
	MongoPoolSize       int    `mapstructure:"mongo-pool-size"      desc:"Maximum size of the connection pool"          default:"4096"`
	MongoTLSCA          string `mapstructure:"mongo-custom-ca"      desc:"Custom certificate authority"`
	MongoTLSCertificate string `mapstructure:"mongo-tls-cert"       desc:"Path to the client certificate"`
	MongoTLSDisable     bool   `mapstructure:"mongo-tls-disable"    desc:"Set this to completely disable TLS"           hidden:"true"`
	MongoTLSSkip        bool   `mapstructure:"mongo-tls-skip"       desc:"Skip CA verification"`
	MongoTLSKey         string `mapstructure:"mongo-tls-key"        desc:"Path to the client key"`
	MongoTLSKeyPass     string `mapstructure:"mongo-tls-key-pass"   desc:"Password for the client key"                  secret:"true" file:"true"`
	MongoURL            string `mapstructure:"mongo-url"            desc:"Mongo connection string"                      required:"true"`
	MongoUser           string `mapstructure:"mongo-user"           desc:"User to use to connect to MongoDB"`
}

// TLSConfig returns the configured TLS configuration as *tls.Config.
func (c *MongoConf) TLSConfig() (*tls.Config, error) {

	if c.MongoTLSDisable {
		return nil, nil
	}

	tlscfg := &tls.Config{}

	if c.MongoTLSCA == "" {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to load system cert pool: %w", err)
		}
		tlscfg.RootCAs = pool
	} else {
		caData, err := os.ReadFile(c.MongoTLSCA)
		if err != nil {
			return nil, fmt.Errorf("unable to load ca file: %w", err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		tlscfg.RootCAs = pool
	}

	if c.MongoTLSCertificate != "" {
		cert, key, err := tglib.ReadCertificatePEM(c.MongoTLSCertificate, c.MongoTLSKey, c.MongoTLSKeyPass)
		if err != nil {
			return nil, fmt.Errorf("unable to load client certificate: %w", err)
		}
		tlscert, err := tglib.ToTLSCertificate(cert, key)
		if err != nil {
			return nil, fmt.Errorf("unable to convert to tls.Certificate: %w", err)
		}
		tlscfg.Certificates = []tls.Certificate{tlscert}
	}

	if c.MongoTLSSkip {
		tlscfg.InsecureSkipVerify = true
	}

	return tlscfg, nil
}

// NATSConf holds the configuration for pubsub connection.
type NATSConf struct {
	NATSClientID       string `mapstructure:"nats-client-id"                 desc:"Nats client ID"`
	NATSClusterID      string `mapstructure:"nats-cluster-id"                desc:"Nats cluster ID"                                default:"test-cluster"`
	NATSPassword       string `mapstructure:"nats-pass"                      desc:"Password to use to connect to Nats"             secret:"true" file:"true"`
	NATSTLSCA          string `mapstructure:"nats-tls-ca"                    desc:"Path to the CA used by Nats"`
	NATSTLSCertificate string `mapstructure:"nats-tls-cert"                  desc:"Path to the client certificate"`
	NATSTLSDisable     bool   `mapstructure:"nats-tls-disable"               desc:"Disable TLS completely"`
	NATSTLSKey         string `mapstructure:"nats-tls-key"                   desc:"Path to the client key"`
	NATSTLSKeyPass     string `mapstructure:"nats-tls-key-pass"              desc:"Password for the client key"                    secret:"true" file:"true"`
	NATSTLSSkip        bool   `mapstructure:"nats-tls-skip"                  desc:"Skip CA verification"`
	NATSURL            string `mapstructure:"nats-url"                       desc:"URL of the nats service"`
	NATSUser           string `mapstructure:"nats-user"                      desc:"User name to use to connect to Nats"            secret:"true" file:"true"`
}

// TLSConfig returns the configured TLS configuration as *tls.Config.
func (c *NATSConf) TLSConfig() (*tls.Config, error) {

	if c.NATSTLSDisable {
		return nil, nil
	}

	tlscfg := &tls.Config{}

	if c.NATSTLSCA == "" {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to load system cert pool: %w", err)
		}
		tlscfg.RootCAs = pool
	} else {
		caData, err := os.ReadFile(c.NATSTLSCA)
		if err != nil {
			return nil, fmt.Errorf("unable to load ca file: %w", err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		tlscfg.RootCAs = pool
	}

	if c.NATSTLSCertificate != "" {
		cert, key, err := tglib.ReadCertificatePEM(c.NATSTLSCertificate, c.NATSTLSKey, c.NATSTLSKeyPass)
		if err != nil {
			return nil, fmt.Errorf("unable to load client certificate: %w", err)
		}
		tlscert, err := tglib.ToTLSCertificate(cert, key)
		if err != nil {
			return nil, fmt.Errorf("unable to convert to tls.Certificate: %w", err)
		}
		tlscfg.Certificates = []tls.Certificate{tlscert}
	}

	if c.NATSTLSSkip {
		tlscfg.InsecureSkipVerify = true
	}

	return tlscfg, nil
}

// NATSPublisherConf holds the config a Pubsub publisher.
type NATSPublisherConf struct {
	NATSPublishTopic string `mapstructure:"nats-publish-topic"        desc:"Topic to use to push events"                 default:"events"`

	NATSConf `mapstructure:",squash"`
}

// NATSConsumerConf holds the config a Pubsub consumer.
type NATSConsumerConf struct {
	NATSGroupName      string `mapstructure:"nats-group-name"           desc:"Nats group name"                             default:"main"`
	NATSSubscribeTopic string `mapstructure:"nats-subscribe-topic"      desc:"Topic to use to receive updates"             default:"override-me"`

	NATSConf `mapstructure:",squash"`
}

// MTLSHeaderConf holds the conf for the secure MTLS header.
type MTLSHeaderConf struct {
	Enabled    bool   `mapstructure:"mtls-header-enabled"    desc:"Trust the value of the defined header containing a user certificate. This is insecure if there is no proper tls verification happening upstream"`
	HeaderKey  string `mapstructure:"mtls-header-key"        desc:"The header to check for user certificates" default:"x-tls-certificate"`
	Passphrase string `mapstructure:"mtls-header-passphrase" desc:"The passphrase to decrypt the AES encrypted header content. It is mandatory if --mtls-header-enabled is set."`
}

// A3SClientConf holds a3s config.
type A3SClientConf struct {
	A3SURL                  string `mapstructure:"a3s-url"               desc:"URL of the a3s server"                                `
	A3SNamespace            string `mapstructure:"a3s-namespace"         desc:"Namespace"`
	A3SCertificateAuthority string `mapstructure:"a3s-cacert"            desc:"Path to the CA certificate"                           secret:"true" file:"true"`
	A3SClientCert           string `mapstructure:"a3s-cert"              desc:"Path to the client certificate"                       secret:"true" file:"true"`
	A3SClientKey            string `mapstructure:"a3s-key"               desc:"Path to the client key"                               secret:"true" file:"true"`
	A3SClientKeyPass        string `mapstructure:"a3s-key-pass"          desc:"Password for the client key"                          secret:"true" file:"true"`
	A3SourceName            string `mapstructure:"a3s-source-name"       desc:"Name of the source to utilize by default"             default:"gateway"`

	systemCAPool *x509.CertPool
}

// SystemCAPool returns the system signing pool
func (c *A3SClientConf) SystemCAPool() (*x509.CertPool, error) {

	if c.A3SCertificateAuthority == "" {
		return nil, fmt.Errorf("no system certificate provided")
	}

	if c.systemCAPool != nil {
		return c.systemCAPool, nil
	}

	data, err := os.ReadFile(c.A3SCertificateAuthority)
	if err != nil {
		return nil, err
	}

	c.systemCAPool = x509.NewCertPool()
	if !c.systemCAPool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("unable to append system signing ca")
	}

	return c.systemCAPool, nil
}

// GatewayConf holds the configuration for the bahamut gateway behaviors.
type GatewayConf struct {
	GWAnnouncedAddress string   `mapstructure:"gw-announce-address" desc:"If set, announce as the service address to the gateway"`
	GWTopic            string   `mapstructure:"gw-topic"            desc:"Topic to use for gateway services discovery"`
	GWOverridePrivate  []string `mapstructure:"gw-override-private" desc:"Overrides the api public/private. In form <name>:<override>. namespace:private makes namespaces api private on the gateway"`
	GWAnnouncePrefix   string   `mapstructure:"gw-announce-prefix"  desc:"Sets the prefix to use for the bahaamut gateway announcement"`
}

// GWPrivateOverrides returns the private overrides in the needed format.
func (c *GatewayConf) GWPrivateOverrides() map[elemental.Identity]bool {

	out := map[elemental.Identity]bool{}

	for _, v := range c.GWOverridePrivate {
		parts := strings.SplitN(v, ":", 2)
		identity := api.Manager().IdentityFromAny(parts[0])

		if parts[0] == "*" {
			for _, i := range api.AllIdentities() {
				out[i] = parts[1] == "public"
			}
			continue
		}
		out[identity] = parts[1] == "public"
	}

	return out
}
