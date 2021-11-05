package token

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Verify verifies the jwt locally using the given certificate.
func Verify(tokenString string, cert *x509.Certificate) (*IdentityToken, error) {

	c := &IdentityToken{}

	_, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			return cert.PublicKey.(*ecdsa.PublicKey), nil
		}
		return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
	})
	if err != nil {
		return nil, fmt.Errorf("unable to parse jwt: %w", err)
	}

	return c, nil
}
