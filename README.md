# a3s

> NOTE: This is a work in progress.

a3s (stands for Auth As A Service) is an authentication and ABAC authorization
server.

It allows to normalize various sources of authentication like OIDC,
AWS/Azure/GCP Identity tokens, LDAP and more into a generic authentication token
that contains identity claims (rather than scopes). These claims can be used by
authorization policies to give a particular subset of users various permissions.

These authorization policies match a set of users based on a logical claim
expression (like `group=red and color=blue or group=admin`) and apply to a
namespace.

A namespace is a node that is part of hierarchical tree that represent an
abstract organizational unit.

Basically, an authorization policy allows a subset of users, defined by claims
retrieved from an authentication source, to perform actions in a particular
namespace and all of its children.

## Quick start

The easiest way to get started is to use the `docker-compose.yaml` in the `dev`
folder.

First, generate the needed certificates:

	dev/certs-init

This creates the needed certificates in `dev/.data/certificates` that the a3s
container will mount (the same certificates will be used by the dev
environment).

Then build the docker container: 

	make docker

And finally start the docker-compose file:

	cd ./dev
	docker compose up

## Using the system

You can start interact with the system by using the raw API with curl or using
the provided `a3sctl`. The later provide a streamlined interface that makes it
more pleasant to use than the raw API.

### Install a3sctl

To install the cli, run:

	make cli

This will install `a3sctl` into you`$GOBIN` folder. You should have this folder
in your `$PATH` if you want to use the cli without using its full path.

### Obtain a root token 

In order to configure the system and create additional namespaces, additional
namespaces, authorizations, etc, you need to obtain a root token to start
interacting with the server:

	 a3sctl auth mtls \
		--api https://127.0.0.1:44443 \
		--api-skip-verify \
		--cert dev/.data/certificates/user-cert.pem \
		--key dev/.data/certificates/user-key.pem \
		--source-name root \

> NOTE: In production environment, never use --api-skip-verify. You should
> instead trust the CA used to issue a3s TLS certificate.

This will print a token you can use for subsequent calls. You can set in the
`$A3SCTL_TOKEN` env variable to use it automatically in the subsequent calls.

If you want to check the content of a token, you can use:

	$ a3sctl auth check --token <token>
	alg: ES256
	kid: 1DAA6949AACB82DBEF1CFE7D93586DD0BF1F090A

	{
	  "exp": 1636830341,
	  "iat": 1636743941,
	  "identity": [
		"commonname=Jean-Michel",
		"serialnumber=219959457279438724775594138274989969558",
		"fingerprint=C8BB0E5FA7644DDC97FD54AEF09053E880EDA939",
		"issuerchain=D98F838F491542CC238275763AA06B7DC949737D",
		"@sourcetype=mtls",
		"@sourcenamespace=/",
		"@sourcename=root"
	  ],
	  "iss": "https://127.0.0.1",
	  "jti": "b2b441a0-5283-4586-baa7-4a45147aaf46"
	}

You can omit `--token` if you have set `$A3SCTL_TOKEN`.

### Test with the example

There is a very small python Flask server located in `/example/python/testapp`.
This comes with a script that you can inspect that will create a namespace to
handle this application, an MTLS source and two authorizations.

You can take a look at the [README](examples/python/testapp/README.md) in that
folder to get started.

## Dev environment

### Prerequesites

First, clone this repository and make sure you have the following installed:

- go
- mongodb
- nats-server
- tmux & tmuxinator

### Initialize the environment

If this is the first time you start the environment, you need to initialize
various things.  

First, initialize the needed certificates:

	dev/certs-init

Then initialize the database:

	dev/mongo-init

All the development data stored stored in `./dev/.data`. If you delete this
folder, you can reinitialize the environment.

All of a3s configuration is defined as env variables from `dev/env`

Finally, you must initialize the root permissions. a3s has no exception or
hardcoded or weak credentials, so we must add an auth source and an
authorization.

To do so, you run:

    dev/a3s-run --init --init-root-ca dev/.data/certificates/ca-acme-cert.pem

### Start everything

Once initialized, start the tmux session by running:

	dev/env-run

This will launch a tmux session starting everything and giving you a working
terminal. To exit:

	env-kill
