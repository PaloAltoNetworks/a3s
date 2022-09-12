# Python Test Server

This folder contains a very simple Flask server that demonstrate how to
integrate a3s within your application. The server handles three routes:

* `/`: public endpoint
* `/secret`: authenticated endpoint
* `/topsecret`: authenticated endpoint

There is a script `init.sh` that will create the needed resources in
a3s to handle this server.

> NOTE: The script assumes a3s is listening on `https://127.0.0.1:44443`.

The script will:

* Create a certificate authority
* Issue a client certificate for John from the CA
* Issue a client certificate for Michael from the CA
* Clean up existing /testapp namespace
* Create a brand new /testapp namespace
* Create a MTLS source in /testapp to recognize the CA
* Add an authorization to allow John to access `/secret`
* Add an authorization to allow Micheal to access `/secret` and `/topsecret`

> NOTE: The authorization policies created are not completely secure, as they
> should use the issuerchain and/or the fingerprint identity claims. But we keep
> things simple here.

## Install requirements

The server depends on `Flask`, `requests` and `pyopenssl` that you must install:

    pip install flask requests pyopenssl

Note your system may call pip, `pip3`

## Launch the script

> NOTE: Everytime you restart the script, the certificates wil be regenerated
> and the namespace /testapp deleted then recreated.

    ./init.sh

The script will output 2 a3sctl commands to retrieve an identity token using
either John or Michael TLS client certificates. You can store these tokens in
`$TOKEN_J` and `$TOKEN_M`.

## Launch the server

To start the server, run:

    ./server.py

Note you may need to call this script as `python3 ./server.py`

## Try the authorizations

You can try the following curl commands:

    curl -k https://127.0.0.1:5000/secret                       # should fail
    curl -k https://127.0.0.1:5000/topsecret                    # should fail
    curl -k -u Bearer:$TOKEN_J https://127.0.0.1:5000/secret    # should work
    curl -k -u Bearer:$TOKEN_J https://127.0.0.1:5000/topsecret # should fail
    curl -k -u Bearer:$TOKEN_M https://127.0.0.1:5000/secret    # should work
    curl -k -u Bearer:$TOKEN_M https://127.0.0.1:5000/topsecret # should work
