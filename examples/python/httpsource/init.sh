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

tg cert issue --is-ca --name "ca"

tg cert issue --auth-server \
	--name "httpsource" \
	--ip 127.0.0.1 \
	--dns localhost \
	--signing-cert "certs/ca-cert.pem" \
	--signing-cert-key "certs/ca-key.pem"

tg cert issue --auth-client \
	--name "access" \
	--signing-cert "certs/ca-cert.pem" \
	--signing-cert-key "certs/ca-key.pem"

export A3SCTL_API="https://127.0.0.1:44443"
export A3SCTL_API_SKIP_VERIFY="true"

echo "* Retrieving an admin token"
A3SCTL_TOKEN="$(
	a3sctl auth mtls \
		--cert ../../../dev/.data/certificates/user-cert.pem \
		--key ../../../dev/.data/certificates/user-key.pem \
		--source-name root --source-namespace /
)"
export A3SCTL_TOKEN

echo
echo "* Deleting existing /httpsource namespace"
a3sctl api delete namespace "/httpsource" -n /

echo
echo "* Creating /httpsource namespace"
a3sctl api create namespace --with.name "httpsource" -n "/" ||
	die "unable to create /httpsource namespace"

echo
echo "* Creating mtlssource"
a3sctl api create httpsource -n "/httpsource" \
	--with.name "default" \
	--with.ca "$(cat certs/ca-cert.pem)" \
	--with.url https://127.0.0.1:5002/login \
	--with.certificate "$(cat certs/access-cert.pem)" \
	--with.key "$(cat certs/access-key.pem)" ||
	die "unable to create mtls resource"

echo
echo "* Success"

echo
echo "Here is a command to check the modified claims:"
echo
cat <<EOF
	a3sctl auth check --token "\$(
		a3sctl auth http \\
		--api $A3SCTL_API \\
		--api-skip-verify \\
		--source-namespace /httpsource \\
		--user john \\
		--pass secret \\
	)"
EOF
