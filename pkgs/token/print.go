package token

import (
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mdp/qrterminal"
)

type printCfg struct {
	raw     bool
	decoded bool
	qrcode  bool
}

// PrintOption represents options that can be passed to token.Print
type PrintOption func(*printCfg)

// PrintOptionRaw sets the printer to
// print the raw token.
func PrintOptionRaw(enabled bool) PrintOption {
	return func(cfg *printCfg) {
		cfg.raw = enabled
	}
}

// PrintOptionDecoded prints the information
// contained in the token.
func PrintOptionDecoded(enabled bool) PrintOption {
	return func(cfg *printCfg) {
		cfg.decoded = enabled
	}
}

// PrintOptionQRCode prints the token as a QRCode.
func PrintOptionQRCode(enabled bool) PrintOption {
	return func(cfg *printCfg) {
		cfg.qrcode = enabled
	}
}

// Fprint prints the given token string using
// the methods passed as options in the given io.Writer.
// If you pass no option, this function is a noop
func Fprint(w io.Writer, token string, opts ...PrintOption) error {

	cfg := printCfg{}
	for _, o := range opts {
		o(&cfg)
	}

	var addLine bool

	if cfg.decoded {
		if err := printDecoded(w, token); err != nil {
			return err
		}
		addLine = true
	}

	if cfg.qrcode {
		if addLine {
			_, err := fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}
		printQRCode(w, token)
		addLine = true
	}

	if cfg.raw {
		if addLine {
			_, err := fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}
		printRaw(w, token)
	}

	return nil
}

func printDecoded(w io.Writer, token string) error {

	claims := jwt.MapClaims{}
	p := jwt.Parser{}

	t, _, err := p.ParseUnverified(token, &claims)
	if err != nil {
		return err
	}

	data, err := prettyjson.Marshal(claims)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "alg:", t.Method.Alg())
	fmt.Fprintln(w, "kid:", t.Header["kid"])
	if exp, ok := claims["exp"].(float64); ok {
		remaining := time.Until(time.Unix(int64(exp), 0))
		if remaining <= 0 {
			fmt.Fprintln(w, "exp: the token has expired", -remaining.Truncate(time.Second), "ago")
		} else {
			fmt.Fprintln(w, "exp:", remaining.Truncate(time.Second))
		}
	}
	fmt.Fprintln(w, string(data))

	return nil
}

func printQRCode(w io.Writer, token string) {

	qrterminal.GenerateWithConfig(
		token,
		qrterminal.Config{
			Writer:         w,
			Level:          qrterminal.M,
			HalfBlocks:     true,
			QuietZone:      1,
			BlackChar:      qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			WhiteChar:      qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
		},
	)
}

func printRaw(w io.Writer, token string) {
	fmt.Fprintln(w, token)
}
