# Claim Modifier Example

This folder contains an example for a remote identity modifier writtent in
Flask. The server handles the `/mod` endpoint that blindly replaces the claims
by `hello=world`

This server is mtls only. you can only access the server using a client
certificate generated from the generated CA.

There is a script called `init.sh` that initializes the needed resources

> NOTE: The script assumes a3s is listening on `https://127.0.0.1:44443`.

The script will:

* Create a certificate authority
* Issue a server certificate for HTTPS
* Issue a client certificate to use when configuring the the source modifier.
* Clean up existing /claimmod namespace
* Create a brand new /claimmod namespace
* Create a MTLS source in /claimmod that uses the server as a modifier.

## Install requirements

The server depends on `Flask` and `requests` that you must install:

    pip install flask requests

## Launch the script

> NOTE: Everytime you restart the script, the certificates wil be regenerated
> and the namespace /claimmod deleted then recreated.

    ./init.sh

## Launch the server

To start the server, run:

    ./server.py

## Check the claims

The init script will print a command you can run to obtain a token from the
configured source using the modifier in the /claimmod namespace. You should be
able to see something like:

    a3sctl auth check --token "$(
      a3sctl auth mtls \
      --api https://127.0.0.1:44443 \
      --api-skip-verify \
      --source-namespace /claimmod \
      --cert certs/john-cert.pem \
      --key certs/john-key.pem
    )"

    alg: ES256
    kid: B47882D62DE6523090D5F3CA4C7E77B746821523DAC7E5F9A61697ECD292BE61

    {
      "aud": [
        "https://127.0.0.1:44443"
      ],
      "exp": 1638487920,
      "iat": 1638401520,
      "identity": [
        "@source:name=default",
        "@source:namespace=/claimmod",
        "@source:type=mtls",
        "hello=world"
      ],
      "iss": "https://127.0.0.1:44443",
      "jti": "853f18fd-7746-4047-b7a4-f22c4acdfada"
    }

<!-- vim:ts=4:sw=3:sts=4:expandtab
-->
