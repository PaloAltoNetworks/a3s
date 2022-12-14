package ui

import (
	"bytes"
	"embed"
)

var (
	apiURLPlaceholder      = []byte("__API_URL__")
	redirectAPIPlaceholder = []byte("__REDIRECT_URL__")
	audiencePlaceholder    = []byte("__AUDIENCE__")
)

//go:embed js/login
var f embed.FS

// GetLogin returns the login page.
func GetLogin(api string, redirect string, audience string) ([]byte, error) {

	doc, err := f.ReadFile("js/login/dist/index.html")
	if err != nil {
		return nil, err
	}

	doc = bytes.Replace(doc, apiURLPlaceholder, []byte(api), 1)
	doc = bytes.Replace(doc, redirectAPIPlaceholder, []byte(redirect), 1)
	doc = bytes.Replace(doc, audiencePlaceholder, []byte(audience), 1)

	return doc, err
}
