package authlib

import (
	"context"
	"time"

	"go.aporeto.io/manipulate"
	"go.uber.org/zap"
)

var tickDuration = 1 * time.Minute

// TokenIssuerFunc is the type of function that can be used
// to retrieve a token.
type TokenIssuerFunc func(context.Context, time.Duration) (string, error)

// A PeriodicTokenManager issues an renew tokens periodically.
type PeriodicTokenManager struct {
	validity   time.Duration
	issuerFunc TokenIssuerFunc
}

// NewPeriodicTokenManager returns a new PeriodicTokenManager backed by midgard.
func NewPeriodicTokenManager(validity time.Duration, issuerFunc TokenIssuerFunc) *PeriodicTokenManager {

	if issuerFunc == nil {
		panic("issuerFunc cannot be nil")
	}

	return &PeriodicTokenManager{
		issuerFunc: issuerFunc,
		validity:   validity,
	}
}

// Issue issues a token.
func (m *PeriodicTokenManager) Issue(ctx context.Context) (token string, err error) {

	return m.issuerFunc(ctx, m.validity)
}

// Run runs the token renewal job.
func (m *PeriodicTokenManager) Run(ctx context.Context, tokenCh chan string) {

	nextRefresh := time.Now().Add(m.validity / 2)

	for {

		select {

		case <-time.After(tickDuration):

			now := time.Now()
			if now.Before(nextRefresh) {
				break
			}

			subctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			token, err := m.Issue(subctx)
			cancel()

			if err != nil {
				zap.L().Error("Unable to renew token", zap.Error(err))
				break
			}

			tokenCh <- token

			nextRefresh = now.Add(m.validity / 2)

		case <-ctx.Done():
			return
		}
	}
}

type x509TokenManager struct {
	m manipulate.Manipulator
	c *Client

	*PeriodicTokenManager
}

func (t *x509TokenManager) SetManipulator(m manipulate.Manipulator) {
	t.m = m
	t.c = NewClient(t.m)
}

// NewX509TokenManager returns a new X.509 backed manipulate.SelfTokenManager.
func NewX509TokenManager(
	sourceName string,
	sourceNamespace string,
	opts ...Option,
) manipulate.SelfTokenManager {

	cfg := newConfig()
	for _, o := range opts {
		o(&cfg)
	}

	t := &x509TokenManager{
		PeriodicTokenManager: &PeriodicTokenManager{
			validity: cfg.validity,
		},
	}

	t.issuerFunc = func(ctx context.Context, v time.Duration) (string, error) {
		return t.c.AuthFromCertificate(
			ctx,
			sourceName,
			sourceNamespace,
			opts...,
		)
	}

	return t
}
