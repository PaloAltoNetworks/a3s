#!/bin/bash

set -e

mkdir -p certs
export TLSGEN_OUT="./certs"
export TLSGEN_FORCE="true"

echo "* Generating certificates..."
tg cert issue --is-ca --name "a3s-test-authority"
tg cert issue --auth-client \
	--name "john" \
	--signing-cert "certs/a3s-test-authority-cert.pem" \
	--signing-cert-key "certs/a3s-test-authority-key.pem"

tg cert issue --auth-client \
	--name "michael" \
	--signing-cert "certs/a3s-test-authority-cert.pem" \
	--signing-cert-key "certs/a3s-test-authority-key.pem"

export A3SCTL_API="https://127.0.0.1:44443"
export A3SCTL_API_SKIP_VERIFY="true"

echo
echo "* Deleting existing /testapp namespace"
a3sctl api delete namespace "/testapp" -n /

echo
echo "* Creating /testapp namespace"
a3sctl api create namespace --name "testapp" -n "/"

echo
echo "* Creating mtlssource"
a3sctl api create mtlssource -n "/testapp" \
	--name "default" \
	--certificateAuthority "$(cat certs/a3s-test-authority-cert.pem)"

echo
echo "* Creating authorization for /secret"
a3sctl api create authorization -n "/testapp" \
	--name "secret-access" \
	--targetNamespace "/testapp" \
	--subject '[
		[
			"@sourcetype=mtls",
			"@sourcename=default",
			"@sourcenamespace=/testapp",
			"commonname=john"
		]
	]' \
	--permissions '["/secret,GET"]'

echo
echo "* Creating authorization for /topsecret"
a3sctl api create authorization -n "/testapp" \
	--name "top-secret-access" \
	--targetNamespace "/testapp" \
	--subject '[
		[
			"@sourcetype=mtls",
			"@sourcename=default",
			"@sourcenamespace=/testapp",
			"commonname=michael"
		]
	]' \
	--permissions '[
		"/secret,GET",
		"/topsecret,GET"
	]'

echo
echo "* Success"

echo
echo "Here is a command to get a token for john:"
echo
cat <<EOF
	a3sctl auth mtls \\
		--audience testapp \\
		--source-namespace /testapp \\
		--cert certs/john-cert.pem \\
		--key certs/john-key.pem
EOF

echo
echo "Here is a command to get a token for michael:"
echo
cat <<EOF
	a3sctl auth mtls \\
		--audience testapp \\
		--source-namespace /testapp \\
		--cert certs/michael-cert.pem \\
		--key certs/michael-key.pem
EOF
echo
