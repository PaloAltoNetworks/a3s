package helpers

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

// ReadFlag will try to read the flag with the given key.
//
// * If it is set to anything but '-, the function will return
// the value as is.
// * If is is set to '-', the function will prompt the user
// to enter the flag using the given title.
// * If secret is true, the function will treat the input
// as a password, not echoing keystrokes.
func ReadFlag(title string, key string, secret bool) string {

	pass := viper.GetString(key)

	if pass != "-" {
		return pass
	}

	if title != "" {
		fmt.Fprint(os.Stderr, title) // nolint: errcheck
	}

	var value []byte
	var err error
	if secret {
		value, err = term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprint(os.Stderr, "\n") // nolint: errcheck
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		value = scanner.Bytes()
	}

	if err != nil {
		panic("unable to read your information: %s")
	}

	return string(value)
}
