#!/bin/bash

set -e

function die() {
	echo "Error: $1"
	exit 1
}

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
echo "* Retrieving an admin token"
A3SCTL_TOKEN="$(
	a3sctl auth mtls \
		--cert ../../../dev/.data/certificates/user-cert.pem \
		--key ../../../dev/.data/certificates/user-key.pem \
		--source-name root --source-namespace /
)"
export A3SCTL_TOKEN

echo
echo "* Deleting existing /testapp namespace"
a3sctl api delete namespace "/testapp" -n /

echo
echo "* Creating /testapp namespace"
a3sctl api create namespace --name "testapp" -n "/" ||
	die "unable to create /testapp namespace"

echo
echo "* Creating mtlssource"
a3sctl api create mtlssource -n "/testapp" \
	--name "default" \
	--certificate-authority "$(cat certs/a3s-test-authority-cert.pem)" ||
	die "unable to create mtls resource"

echo
echo "* Creating authorization for /secret"
a3sctl api create authorization -n "/testapp" \
	--name "secret-access" \
	--target-namespace "/testapp" \
	--subject '[
		[
			"@sourcetype=mtls",
			"@sourcename=default",
			"@sourcenamespace=/testapp",
			"commonname=john"
		]
	]' \
	--permissions '["/secret,GET"]' ||
	die "unable to create authorization for /secret"

echo
echo "* Creating authorization for /topsecret"
a3sctl api create authorization -n "/testapp" \
	--name "top-secret-access" \
	--target-namespace "/testapp" \
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
	]' ||
	die "unable to create authorization for /topsecret"

echo
echo "* Success"

echo
echo "Here is a command to get a token for john:"
echo
cat <<EOF
	a3sctl auth mtls \\
		--api https://127.0.0.1:44443 \\
		--api-skip-verify \\
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
		--api https://127.0.0.1:44443 \\
		--api-skip-verify \\
		--audience testapp \\
		--source-namespace /testapp \\
		--cert certs/michael-cert.pem \\
		--key certs/michael-key.pem
EOF
echo
