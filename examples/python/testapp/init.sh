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
echo "* Deleting/Creating existing /testapp namespace"
a3sctl api delete namespace "/testapp" -n /
a3sctl api create namespace --with.name "testapp" -n "/" ||
	die "unable to create /testapp namespace"

echo
echo "Importing data"
a3sctl api create import \
	-n /testapp \
	--input-file=import.gotmpl ||
	die "unable to import data"

echo
echo "* Success"

echo
echo "Here is a command to get a token for john:"
echo
cat <<EOF
	export jtok=\$( \\
		a3sctl auth mtls \\
			--api $A3SCTL_API \\
			--api-skip-verify \\
			--audience testapp \\
			--source-namespace /testapp \\
			--cert certs/john-cert.pem \\
			--key certs/john-key.pem \\
	)
EOF

echo
echo "Here is a command to get a token for michael:"
echo
cat <<EOF
	export mtok=\$( \\
		a3sctl auth mtls \\
			--api $A3SCTL_API \\
			--api-skip-verify \\
			--audience testapp \\
			--source-namespace /testapp \\
			--cert certs/michael-cert.pem \\
			--key certs/michael-key.pem \\
	)
EOF
echo
