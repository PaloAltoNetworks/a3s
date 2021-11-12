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

## Dev environment

### Prerequesites

First, clone this repository and make sure you have the following installed:

- go
- mongod
- nats-server
- tmux & tmuxinator

### Initialize environment

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

Once initialized, go to the `dev` folder and run:

	./env-run

This will launch a tmux session starting everything and giving you a working
terminal. To exit:

	env-kill
