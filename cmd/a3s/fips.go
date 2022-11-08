//go:build boringcrypto

package main

import (
	"crypto/boring"
	"fmt"
)

func init() {
	fmt.Println("FIPS: boringcrypto enabled:", boring.Enabled())
}
