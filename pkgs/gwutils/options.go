package gwutils

import (
	"time"
)

// VerifierOption can be used to configure
// optional aspect of MakeTLSPeerCertificateVerifier.
type VerifierOption func(*verifierConf)

type verifierConf struct {
	cacheDuration time.Duration
	cacheMaxSize  int64
	timeout       time.Duration
}

func newVerifierConf() verifierConf {
	return verifierConf{
		cacheMaxSize:  2048,
		cacheDuration: 10 * time.Minute,
		timeout:       10 * time.Second,
	}
}

// OptionCacheDuration sets the life time of cached CAs.
func OptionCacheDuration(d time.Duration) VerifierOption {
	return func(cfg *verifierConf) {
		cfg.cacheDuration = d
	}
}

// OptionCacheSize sets the maximum number of items
// that can be in the cache, before evicting older ones.
func OptionCacheSize(s int64) VerifierOption {
	return func(cfg *verifierConf) {
		cfg.cacheMaxSize = s
	}
}

// OptionTimeout sets the maximum amount of time to
// wait for A3S to reply.
func OptionTimeout(d time.Duration) VerifierOption {
	return func(cfg *verifierConf) {
		cfg.timeout = d
	}
}
