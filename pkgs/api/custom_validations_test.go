package api

import (
	"fmt"
	"testing"
)

func TestValidateCIDR(t *testing.T) {
	type args struct {
		attribute string
		network   string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"valid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"10.0.1.0/24",
				}
			},
			false,
			nil,
		},
		{
			"invalid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"10.0.1.024",
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must be a CIDR"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
		{
			"empty cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"",
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must be a CIDR"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := ValidateCIDR(tArgs.attribute, tArgs.network)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCIDR error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestValidateCIDROptional(t *testing.T) {
	type args struct {
		attribute string
		network   string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"valid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"10.0.1.0/24",
				}
			},
			false,
			nil,
		},
		{
			"invalid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"10.0.1.024",
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must be a CIDR"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
		{
			"empty cidr",
			func(*testing.T) args {
				return args{
					"attr",
					"",
				}
			},
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := ValidateCIDROptional(tArgs.attribute, tArgs.network)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCIDROptional error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestValidateCIDRList(t *testing.T) {
	type args struct {
		attribute string
		networks  []string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"valid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					[]string{"10.0.1.0/24", "11.0.1.0/24"},
				}
			},
			false,
			nil,
		},
		{
			"invalid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					[]string{"10.0.1.0/24", "11.0.1.024"},
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must be a CIDR"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
		{
			"empty cidr",
			func(*testing.T) args {
				return args{
					"attr",
					nil,
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must not be empty"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := ValidateCIDRList(tArgs.attribute, tArgs.networks)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCIDRList error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestValidateCIDRListOptional(t *testing.T) {
	type args struct {
		attribute string
		networks  []string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"valid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					[]string{"10.0.1.0/24", "11.0.1.0/24"},
				}
			},
			false,
			nil,
		},
		{
			"invalid cidr",
			func(*testing.T) args {
				return args{
					"attr",
					[]string{"10.0.1.0/24", "11.0.1.024"},
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: Attribute 'attr' must be a CIDR"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err)
					t.Fail()
				}
			},
		},
		{
			"empty cidr",
			func(*testing.T) args {
				return args{
					"attr",
					nil,
				}
			},
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := ValidateCIDRListOptional(tArgs.attribute, tArgs.networks)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCIDRListOptional error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestValidateTagsExpression(t *testing.T) {
	type args struct {
		attribute  string
		expression [][]string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"empty tag expression",
			func(*testing.T) args {
				return args{
					"attr",
					[][]string{},
				}
			},
			false,
			nil,
		},
		{
			"half empty tag expression",
			func(*testing.T) args {
				return args{
					"attr",
					[][]string{nil, nil},
				}
			},
			false,
			nil,
		},
		{
			"nil tag expression",
			func(*testing.T) args {
				return args{
					"attr",
					nil,
				}
			},
			false,
			nil,
		},
		{
			"valid tag expression",
			func(*testing.T) args {
				return args{
					"attr",
					[][]string{{"a=a", "b=b"}, {"c=c"}},
				}
			},
			false,
			nil,
		},
		{
			"too long tag expression",
			func(*testing.T) args {
				long := make([]byte, 1025)
				return args{
					"attr",
					[][]string{{string(long), "b=b"}, {"c=c"}},
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := fmt.Sprintf("error 422 (a3s): Validation Error: '%s' must be less than 1024 bytes", make([]byte, 1025))
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err.Error())
					t.Fail()
				}
			},
		},
		{
			"invalid tag expression",
			func(*testing.T) args {
				return args{
					"attr",
					[][]string{{"aa", "b=b"}, {"c=c"}},
				}
			},
			true,
			func(err error, t *testing.T) {
				wanted := "error 422 (a3s): Validation Error: 'aa' must contain at least one '=' symbol separating two valid words"
				if err.Error() != wanted {
					t.Logf("wanted %s but got %s", wanted, err.Error())
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := ValidateTagsExpression(tArgs.attribute, tArgs.expression)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTagsExpression error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestValidateAuthorizationSubject(t *testing.T) {
	type args struct {
		attribute string
		subject   [][]string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrString string
	}{
		{
			"valid subject",
			args{
				"subject",
				[][]string{
					{"@auth:realm=certificate", "@auth:claim=a"},
					{"@auth:realm=vince", "@auth:claim=a", "@auth:claim=b"},
				},
			},
			false,
			"",
		},
		{
			"missing realm claim",
			args{
				"subject",
				[][]string{
					{"@auth:realm=certificate", "@auth:claim=a"},
					{"@auth:claim=a", "@auth:claim=b"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject line 2 must contain the '@auth:realm' key",
		},
		{
			"2 realm claims",
			args{
				"subject",
				[][]string{
					{"@auth:realm=certificate", "@auth:claim=a", "@auth:realm=vince"},
					{"@auth:claim=a", "@auth:claim=b"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject line 1 must contain only one '@auth:realm' key",
		},
		{
			"single claim line",
			args{
				"subject",
				[][]string{
					{"@auth:realm=certificate", "@auth:claim=a"},
					{"@auth:realm=certificate"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject and line should contain at least 2 claims",
		},
		{
			"missing auth prefix claim",
			args{
				"subject",
				[][]string{
					{"@auth:realm=certificate", "@auth:claim=a"},
					{"@auth:claim=a", "@auth:claim=b", "not:good"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject claims 'not:good' on line 2 must be prefixed by '@auth:'",
		},
		{
			"oidc correct",
			args{
				"subject",
				[][]string{
					{"@auth:realm=oidc", "@auth:claim=a", "@auth:namespace=/a/b"},
					{"@auth:realm=vince", "@auth:claim=a", "@auth:claim=b"},
				},
			},
			false,
			"",
		},
		{
			"oidc missing namespace",
			args{
				"subject",
				[][]string{
					{"@auth:realm=oidc", "@auth:claim=a"},
					{"@auth:realm=vince", "@auth:claim=a", "@auth:claim=b"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: The realm OIDC mandates to add the '@auth:namespace' key to prevent potential security side effects",
		},
		{
			"saml correct",
			args{
				"subject",
				[][]string{
					{"@auth:realm=saml", "@auth:claim=a", "@auth:namespace=/a/b"},
					{"@auth:realm=vince", "@auth:claim=a", "@auth:claim=b"},
				},
			},
			false,
			"",
		},
		{
			"saml missing namespace",
			args{
				"subject",
				[][]string{
					{"@auth:realm=saml", "@auth:claim=a"},
					{"@auth:realm=vince", "@auth:claim=a", "@auth:claim=b"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: The realm SAML mandates to add the '@auth:namespace' key to prevent potential security side effects",
		},
		{
			"broken tag with no equal",
			args{
				"subject",
				[][]string{
					{"@auth:realm=saml", "@auth:claim"},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject claims '@auth:claim' on line 1 is an invalid tag",
		},
		{
			"broken tag with no value",
			args{
				"subject",
				[][]string{
					{"@auth:realm=saml", "@auth:claim="},
				},
			},
			true,
			"error 422 (a3s): Validation Error: Subject claims '@auth:claim=' on line 1 has no value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthorizationSubject(tt.args.attribute, tt.args.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIAuthorizationPolicySubject() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && err.Error() != tt.wantErrString {
				t.Errorf("ValidateAPIAuthorizationPolicySubject() error = '%v', wantErrString = '%v'", err, tt.wantErrString)
			}
		})
	}
}

func TestValidatePEM(t *testing.T) {
	type args struct {
		attribute string
		pemdata   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"nothing set",
			args{
				"pem",
				``,
			},
			false,
		},
		{
			"valid single PEM",
			args{
				"pem",
				`-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoTCXNlcGhpcm90aDEUMBIGA1UEAxMLYXV0b21hdGlvbnMw
WTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASxKA9vbyk7FXXlOCi0kTKLVne/mK8o
ZQDPRcehze0EMwTAR5loNahC19hQtExCi64fmI3QCcrEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----`,
			},
			false,
		},
		{
			"valid single PEM",
			args{
				"pem",
				`-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoTCXNlcGhpcm90aDEUMBIGA1UEAxMLYXV0b21hdGlvbnMw
WTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASxKA9vbyk7FXXlOCi0kTKLVne/mK8o
ZQDPRcehze0EMwTAR5loNahC19hQtExCi64fmI3QCcrEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoTCXNlcGhpcm90aDEUMBIGA1UEAxMLYXV0b21hdGlvbnMw
WTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASxKA9vbyk7FXXlOCi0kTKLVne/mK8o
ZQDPRcehze0EMwTAR5loNahC19hQtExCi64fmI3QCcrEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----
`,
			},
			false,
		},
		{
			"invalid single PEM",
			args{
				"pem",
				`-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoT ----NOT PEM---- I3QCcrEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----`,
			},
			true,
		},
		{
			"valid single PEM",
			args{
				"pem",
				`-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoTCXNlcGhpcm90aDEUMBIGA1UEAxMLYXV0b21hdGlvbnMw
WTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASxKA9vbyk7FXXlOCi0kTKLVne/mK8o
ZQDPRcehze0EMwTAR5loNahC19hQtExCi64fmI3QCcrEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIBpDCCAUmgAwIBAgIQDbXKAZzk9RjcNSGMsWke1zAKBggqhkjOPQQDAjBGMRAw
DgYDVQQKEwdBcG9yZXRvMQ8wDQYDVQQLEwZhcG9tdXgxITAfBgNVBAMTGEFwb211
eCBQdWJsaWMgU2lnbmluZyBDQTAeFw0xOTAxMjQyMjQ3MjlaFw0yODEyMDIyMjQ3
MjlaMCoxEjAQBgNVBAoTCXNlcGhpcm90aDEUMBIGA1UEAxMLYXV0b21hdGlvbnMw
WTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASxKA9vbyk7FXXlOCi0kTKLVne/mK8o
ZQDPRcehze0EMwTAR5     ----NOT PEM----   crEGH9ycUoITYPgozUwMzAO
BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAKBggqhkjOPQQDAgNJADBGAiEAm1u2T1vRooIy3rd0BmBSAa6WR6BtHl9nDbGN
1ZM+SgsCIQDu4R6OziiWbRdn50bneZT5qPO+07ALY5m4DG96VyCaQw==
-----END CERTIFICATE-----
`,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePEM(tt.args.attribute, tt.args.pemdata); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePEM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
