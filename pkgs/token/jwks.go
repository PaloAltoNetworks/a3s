package token

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

// Various errors returned by a JWKS.
var (
	ErrJWKSNotFound    = errors.New("kid not found in JWKS")
	ErrJWKSInvalidType = errors.New("certificate must be ecdsa")
	ErrJWKSKeyExists   = errors.New("key with the same kid already exists")
)

// A JWKS is a structure to manage a JSON Web Key Set.
type JWKS struct {
	Keys []*JWKSKey `json:"keys"`

	keyMap map[string]*JWKSKey

	sync.RWMutex
}

// NewJWKS returns a new JWKS.
func NewJWKS() *JWKS {
	return &JWKS{
		keyMap: map[string]*JWKSKey{},
	}
}

// Append appends a new certificate to the JWKS.
func (j *JWKS) Append(cert *x509.Certificate) error {
	return j.AppendWithPrivate(cert, nil)
}

// AppendWithPrivate appends a new certificate and its private key to the JWKS.
func (j *JWKS) AppendWithPrivate(cert *x509.Certificate, private crypto.PrivateKey) error {

	j.Lock()
	defer j.Unlock()

	public, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return ErrJWKSInvalidType
	}

	kid := fmt.Sprintf("%02X", sha1.Sum(cert.Raw))

	if _, ok := j.keyMap[kid]; ok {
		return ErrJWKSKeyExists
	}

	k := &JWKSKey{
		KTY:     "EC",
		KID:     kid,
		Use:     "sign",
		CRV:     public.Curve.Params().Name,
		X:       base64.RawURLEncoding.EncodeToString(public.X.Bytes()),
		x:       public.X,
		Y:       base64.RawURLEncoding.EncodeToString(public.Y.Bytes()),
		y:       public.Y,
		private: private,
	}

	j.Keys = append(j.Keys, k)
	j.keyMap[kid] = k

	return nil
}

// Get returns the key with the given ID.
// Returns ErrJWKSNotFound if not found.
func (j *JWKS) Get(kid string) (*JWKSKey, error) {

	j.RLock()
	defer j.RUnlock()

	k, ok := j.keyMap[kid]
	if !ok {
		return nil, ErrJWKSNotFound
	}

	return k, nil
}

// GetLast returns the last inserted key.
func (j *JWKS) GetLast() *JWKSKey {

	j.RLock()
	defer j.RUnlock()

	if len(j.Keys) == 0 {
		return nil
	}

	return j.Keys[len(j.Keys)-1]
}

// Del deletes the key with the given ID.
// Returns true if something was deleted, false
// otherwise.
func (j *JWKS) Del(kid string) bool {

	j.Lock()
	defer j.Unlock()

	if _, ok := j.keyMap[kid]; !ok {
		return false
	}

	delete(j.keyMap, kid)

	var idx int
	for i, key := range j.Keys {
		if key.KID == kid {
			idx = i
			break
		}
	}

	j.Keys = append(j.Keys[:idx], j.Keys[idx+1:]...)

	return true
}

// JWKSKey represents a single key stored in
// a JWKS.
type JWKSKey struct {
	KTY string `json:"kty"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg,omitempty"`
	N   string `json:"n,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
	CRV string `json:"crv,omitempty"`

	x       *big.Int
	y       *big.Int
	private crypto.PrivateKey
}

// Curve returns the curve used by the key.
func (k *JWKSKey) Curve() elliptic.Curve {

	switch k.CRV {
	case "P-224":
		return elliptic.P224()
	case "P-256":
		return elliptic.P256()
	case "P-384":
		return elliptic.P384()
	case "P-521":
		return elliptic.P521()
	default:
		return nil
	}
}

// PublicKey returns a ready to use crypto.PublicKey.
func (k *JWKSKey) PublicKey() crypto.PublicKey {

	switch k.KTY {
	case "EC":
		return &ecdsa.PublicKey{
			X:     k.x,
			Y:     k.y,
			Curve: k.Curve(),
		}
	default:
		return nil
	}
}

// PrivateKey returns the crypto.PrivateKey associated to
// the public key, if it was given during adding it.
func (k *JWKSKey) PrivateKey() crypto.PrivateKey {
	return k.private
}
