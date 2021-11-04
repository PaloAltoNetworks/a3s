# a3s

> NOTE: this is a work in progress and this software is not usable yet

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


## Launch Dev environment

> NOTE: This is a work in progress.

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
hardscoded credentials, so we must add an auth source and an authorization. To
do so, you can run:

    dev/a3s-run --init --init-root-ca dev/.data/certificates/ca-acme-cert.pem

### Start everything

Once initialized, go to the `dev` folder and run:

	./env-run

This will launch a tmux session starting everything and giving you a working
terminal. To exit:

	env-kill
