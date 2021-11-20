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

## Table of contents

<!-- vim-markdown-toc GFM -->

* [Quick start](#quick-start)
* [Using the system](#using-the-system)
	* [Install a3sctl](#install-a3sctl)
	* [Obtain a root token](#obtain-a-root-token)
	* [Test with the example](#test-with-the-example)
* [Obtaining identity tokens](#obtaining-identity-tokens)
	* [MTLS](#mtls)
		* [Create an MTLS source](#create-an-mtls-source)
		* [Obtain a token](#obtain-a-token)
	* [LDAP](#ldap)
		* [Create an LDAP source](#create-an-ldap-source)
		* [Obtain a token](#obtain-a-token-1)
	* [OIDC](#oidc)
		* [Create an OIDC source](#create-an-oidc-source)
		* [Obtain a token](#obtain-a-token-2)
	* [Amazon STS](#amazon-sts)
	* [Google Cloud Platform token](#google-cloud-platform-token)
	* [Azure token](#azure-token)
	* [Existing A3S Identity token](#existing-a3s-identity-token)
* [Dev environment](#dev-environment)
	* [Prerequesites](#prerequesites)
	* [Initialize the environment](#initialize-the-environment)
	* [Start everything](#start-everything)
* [Support](#support)
* [Contributing](#contributing)

<!-- vim-markdown-toc -->

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

## Obtaining identity tokens

This section describes how to use the various sources of authentication. and how
to retrieve a token from them.

All these example will assume to work in the namespace `/tutorial`. To create
it, you can run:

	a3sctl api create namespace --name tutorial --namespace /
	export A3SCTL_NAMESPACE=/tutorial

> NOTE: the env variable will tell a3sctl which namespace to target without
> having to pass the `--namespace` flag every time.

> NOTE: Some auth commands will require to pass the namespace of the auth
> source. You can either set `--source-namespace` or leave it empty to fallback on
> the value set by `--namespace`.

> NOTE: you can also get more info about a ressource by using the `-h` flag. This
> will list all the possible properies the api supports

Whichever authentication source you are using, you can always ask for a restricted
token. A restricted token contains additional user requested restrictions
preventing actions that would normally be possible to do based on the
euthorizations associated to the claims.

* `--restrict-namespace`: a namespace restricted token will only be valid if used on the
    restricted namespace or one of its children
* `--restrict-network`: a network restricted token can only be used if the
    source network from which it is used is contained in one of the restricted networks.
* `--restrict-permissions`: limits what permissions the token will have. For instance if your
    authorization set grants `dog:eat,sleep`, you may ask for a token that will only work
    for `dog:eat`.

It is also possible to limit the amount of identity claims that will be emebeded
into the identity token by using the `--cloak` flag. This can be useful for
privacy. For instance if a party request you to have `color=blue` and that is
the only claim that matters, you can hide the rest of your claims by passing

    --cloak color=blue

Cloaking uses prefix matching. So you can decide to only embbed the color
and size claims (if you have multiple of them) by doing:

    --cloak color= --cloak size=

### MTLS

The MTLS source uses mutual TLS to authenticate a client.  The client must
present a client certificate (usage set to auth client) that is signed by the CA
provided in the designed MTLS auth source.

#### Create an MTLS source

You first need to have a CA that can issue certificates for your user. In this
example, we use `tg`, but you can use any PKI tool you like.

	tg cert --name myca --is-ca 
	tg cert --name user1 \
		--signing-cert myca-cert.pem \ 
		--signing-cert-key myca-key.pem

NOTE: Not protecting a private key with a passphrase is bad. Don't do this in
production.

Then we need to create the MTLS auth source:

	a3sctl api create mtlssource \
		--name my-mtls-source \
		--certificate-auhority "$(cat myca-cert.pem)"

#### Obtain a token

To obtain a token from the newly created source:

	a3sctl auth mtls \
		--source-name my-mtls-source \
		--source-namespace /tutorial \
		--cert user1-cert.pem \
		--key user1-key.pem

### LDAP

A3s supports using a remote LDAP as authentication source. The LDAP server must
be accessible from a3s. A3s will refuse to connect to an LDAP with no form of
encryption (TLS or STARTTLS).

#### Create an LDAP source

To create an LDAP source, run:

	a3sctl api create ldapsource \
		--name my-ldap-source \
		--address 127.0.0.1:389 \
		--certificate-authority "$(cat ldap-ce-cert.pem)" \
		--base-dn dc=universe,dc=io \
		--bind-dn cn=readonly,dc=universe,dc=io \
		--bind-password password

* The `base-dn` is the DN to use to search for users.
* Yhe `bind-dn` is the account a3s will use to connect to the ldap. It should be
	a readonly account.
* The `bind-password` is the password associated to the `bind-dn`.

You can also use `--certificate-auhority` to pass a custom CA if the
certificates used by the server are not trusted by the host running a3s.

#### Obtain a token

To obtain a token from the newly created source:

	a3sctl auth ldap \
		--source-name my-ldap-source \
		--namespace /tutorial \
		--user bob \
		--pass s3cr3t

### OIDC

A3s can retrieve an identity token from an existing OIDC provider in order to
deliver normalized identiy tokens.

#### Create an OIDC source

Configuring a valid OIDC provider is beyond the scope of this document. However,
they will all work the same and will give you a client ID, a client secret and
an endpoint.

It is however important to set the allowed redirect URL to be
`http://localhost:65333` on your provider if you plan to use a3sctl to
authenticate.

Once you have this information, create an OIDC source:

	a3sctl api create oidcsource \
		--name my-oidc-source \
		--client-id <client id> \
		--client-secret <client secret> \
		--endpoint https://accounts.google.com \
		--scopes '["email", "given_name"]'

The scopes indicate the OIDC provider which claim to return. This will vary
depending on your provider.

You can also use `--certificate-auhority` to pass a custom CA if the
certificates used by the server are not trusted by the host running a3s.

#### Obtain a token

While all the other sources can be used easily with curl for instance, the OIDC
source necessitate to run a http server and needs to perform a dance that is
painful to do manually. A3sctl will do all of this transparently.

To obtain a token from the newly created source:

	a3sctl auth oidc \
		--source-name my-oidc-source \
		--source-namespace /tutorial

This will print an URL to open in your browser to authenticate against the OIDC
provider. Once done, the provider will call back a3sctl and the token will be
displayed.


### Amazon STS

This authentication source does not need custom source creation as it uses AWS
broadly. How to retrieve a token from AWS is beyond the scope of this document.
However, if you run a3sctl from an EC2 instance that has an IAM role assigned, it
will retrieve one for you, if you don't pass any additional information

If you are not running the command on AWS:

	a3sctl auth aws \
		--access-key-id <kid> \
		--access-key-secret <secret> \
		--access-token <token>

However, if you are running it from an AWS EC2 instance, you just need to run:

	a3sctl auth aws

### Google Cloud Platform token

This authentication source does not need custom source creation as it uses GCP
broadly. How to retrieve a token from GCP is beyond the scope of this document.
However, if you run a3sctl from a GCP instance, it will retrieve one for you, if
you don't pass any additional information

If you are not running the command on GCP:

	a3sctl auth gcp --access-token <token>

However, if you are running it from an GCP instance, you just need to run:

	a3sctl auth gcp

###  Azure token

This authentication source does not need custom source creation as it uses Azure
broadly. How to retrieve a token from Azure is beyond the scope of this document.
However, if you run a3sctl from a Azure instance, it will retrieve one for you, if
you don't pass any additional information

If you are not running the command on GCP:

	a3sctl auth azure --access-token <token>

However, if you are running it from an Azure instance, you just need to run:

	a3sctl auth azure

### Existing A3S Identity token

You can use an existing a3s identity token to ask for another one. Note that is
not a renew mechanism. The requested token cannot expire later than the original
one. The goal of this authentication source is to ask for a more restricted
and/or cloaked version of the original. 

This authentication souurce does not need custom source creation.

To get obtain a token:

    a3sctl auth a3s --token <token> \
        --restrict-namespace /a/child/ns \
        --restrict-network 10.0.1.1/32 \
        --restrict-permissions "dog:eat,sleep" 

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

## Support

Please read [SUPPORT.md](SUPPORT.md) for details on how to get support for this
project.

## Contributing

We value your contributions! Please read
[CONTRIBUTING.md](https://github.com/PaloAltoNetworks/.github/CONTRIBUTING.md)
for details on how to contribute, and the process for submitting pull requests
to us. 
