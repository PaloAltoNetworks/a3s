package token

import (
	"crypto"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
)

func makeToken(claims jwt.Claims, signMethod jwt.SigningMethod, key crypto.PrivateKey) string {

	token := jwt.NewWithClaims(signMethod, claims)
	t, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}

	return t
}

func TestVerifyToken(t *testing.T) {

	cert, key := getECCert()
	cert2, _ := getECCert()

	Convey("Given I verify a valid token", t, func() {

		token := makeToken(
			&jwt.StandardClaims{Subject: "sub"},
			jwt.SigningMethodES256,
			key,
		)

		claims, err := Verify(token, cert)
		So(err, ShouldBeNil)
		So(claims.Valid(), ShouldBeNil)
	})

	Convey("Given I verify a valid token with wrong certificate", t, func() {

		token := makeToken(
			&jwt.StandardClaims{Subject: "sub"},
			jwt.SigningMethodES256,
			key,
		)

		claims, err := Verify(token, cert2)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to parse jwt: crypto/ecdsa: verification error")
		So(claims, ShouldBeNil)
	})
	Convey("Given I verify a valid token with wrong certificate", t, func() {

		token := `eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.kRR_IkmuxhMP56LnBdfSxfKyGVW6DTf6iUChg-L1LHTKyTqtXISJI4PxNJhR40JZ`
		claims, err := Verify(token, cert2)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to parse jwt: unexpected signing method: HS384")
		So(claims, ShouldBeNil)
	})
}
