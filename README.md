# a3s

> NOTE: This is a work in progress.

a3s (stands for Auth As A Service) is an authentication and ABAC authorization
server.

It allows to normalize various sources of authentication like OIDC,
AWS/Azure/GCP tokens, LDAP and more into a generic identity token
that contains identity claims (rather than scopes). These claims can be used by
some authorization policies to give a particular subset of bearers various
permissions.

These authorization policies match a set of bearers based on a logical claim
expression (like `group=red and color=blue or group=admin`) and they apply to a
namespace.

A namespace is a node that is part of hierarchical tree that represent an
abstract organizational unit. The root namespace is named `/`.

Basically, an authorization policy allows a subset of users, defined by claims
retrieved from an authentication source, to perform actions in a particular
namespace and all of its children.

Apps can receive a request alongside a delivered identity token then check with
a3s if the current bearer is allowed to perform a particular action in a
particular namespace.

![flowchart](docs/imgs/Diagram2.png)

## Table of contents

<!-- vim-markdown-toc GFM -->

* [Quick start](#quick-start)
* [Using the system](#using-the-system)
	* [Install a3sctl](#install-a3sctl)
	* [Obtain a root token](#obtain-a-root-token)
	* [Test with the sample app](#test-with-the-sample-app)
* [Obtaining identity tokens](#obtaining-identity-tokens)
	* [Restrictions](#restrictions)
	* [Cloaking](#cloaking)
	* [Autentication sources](#autentication-sources)
		* [MTLS](#mtls)
			* [Create an MTLS source](#create-an-mtls-source)
			* [Obtain a token](#obtain-a-token)
		* [LDAP](#ldap)
			* [Create an LDAP source](#create-an-ldap-source)
			* [Obtain a token](#obtain-a-token-1)
		* [OIDC](#oidc)
			* [Create an OIDC source](#create-an-oidc-source)
			* [Obtain a token](#obtain-a-token-2)
		* [A3S remote identity token](#a3s-remote-identity-token)
			* [Create an A3S source](#create-an-a3s-source)
			* [Obtain a token](#obtain-a-token-3)
		* [Amazon STS](#amazon-sts)
		* [Google Cloud Platform token](#google-cloud-platform-token)
		* [Azure token](#azure-token)
		* [A3S local identity token](#a3s-local-identity-token)
* [Writing authorizations](#writing-authorizations)
	* [Subject](#subject)
	* [Permissions](#permissions)
	* [Target namespaces](#target-namespaces)
	* [Examples](#examples)
* [Check for permissions from your app](#check-for-permissions-from-your-app)
* [Using a3sctl](#using-a3sctl)
	* [Completion](#completion)
		* [Bash](#bash)
		* [Zsh](#zsh)
		* [Fish](#fish)
	* [Configuration file](#configuration-file)
	* [Auto authentication](#auto-authentication)
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
		"@source:name=root"
		"@source:namespace=/",
		"@source:type=mtls",
		"commonname=Jean-Michel",
		"fingerprint=C8BB0E5FA7644DDC97FD54AEF09053E880EDA939",
		"issuerchain=D98F838F491542CC238275763AA06B7DC949737D",
		"serialnumber=219959457279438724775594138274989969558",
	  ],
	  "iss": "https://127.0.0.1",
	  "jti": "b2b441a0-5283-4586-baa7-4a45147aaf46"
	}

You can omit `--token` if you have set `$A3SCTL_TOKEN`.

### Test with the sample app

There is a very small python Flask server located in `/example/python/testapp`.
This comes with a script that you can inspect that will create a namespace to
handle this application: an MTLS source and two authorizations.

You can take a look at the [README](examples/python/testapp/README.md) in that
folder to get started.

## Obtaining identity tokens

This section describes how to use the various sources of authentication and how
to retrieve from them and apply restrictions or cloaking on it.

All following examples will assume to work in the namespace `/tutorial`. To create
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

### Restrictions

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

### Cloaking

It is also possible to limit the amount of identity claims that will be emebeded
into the identity token by using the `--cloak` flag. This can be useful for
privacy reasons. For instance if a party request you to have `color=blue` and that is
the only claim that matters, you can hide the rest of your claims by passing

    --cloak color=blue

Cloaking uses prefix matching. So you can decide to only embbed the color
and size claims (if you have multiple of them) by doing:

    --cloak color= --cloak size=

### Autentication sources

While a3s allows to verify the identity of a token bearer, it does not provide a
way to store information about the users. In order to derive identity claims,
a3s relies on third-party authentication sources, who hold the actual data about
a bearer.

#### MTLS

The MTLS source uses mutual TLS to authenticate a client.  The client must
present a client certificate (usage set to auth client) that is signed by the CA
provided in the designed MTLS auth source.

##### Create an MTLS source

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

##### Obtain a token

To obtain a token from the newly created source:

	a3sctl auth mtls \
		--source-name my-mtls-source \
		--source-namespace /tutorial \
		--cert user1-cert.pem \
		--key user1-key.pem

> NOTE: you can set `-` for '--pass`. In that case, a3sctl will ask for user
> input from stdin.

#### LDAP

A3s supports using a remote LDAP as authentication source. The LDAP server must
be accessible from a3s. A3s will refuse to connect to an LDAP with no form of
encryption (TLS or STARTTLS).

##### Create an LDAP source

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

##### Obtain a token

To obtain a token from the newly created source:

	a3sctl auth ldap \
		--source-name my-ldap-source \
		--namespace /tutorial \
		--user bob \
		--pass s3cr3t

> NOTE: you can set `-` for '--user` and/or `--pass`. In that case, a3sctl will
> ask for user input from stdin.

#### OIDC

A3s can retrieve an identity token from an existing OIDC provider in order to
deliver normalized identiy tokens.

##### Create an OIDC source

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

##### Obtain a token

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

#### A3S remote identity token

This authentication source allows to issue a token from another one issued by
another a3s server. This allows to trust other a3s instances and issue local
tokens from trusted ones.

##### Create an A3S source

You need to create an a3s source in order to validate the remote tokens. This
source requires to pass the raw address of the remote a3s server, as it will use
the well-known jwks URL to retrieve the keys and verify the token signature.

To create an a3s source:

	a3sctl api create a3ssource \
		--name my-remote-a3s-source \
		--issuer https://remote-a3s.com

You can also use `--certificate-auhority` to pass a custom CA if the
certificates used by the server are not trusted by the host running a3s.

If the issuer is not the root URL of the remote a3s server, you can use the
`--endpoint` flag to pass the actual URL.

##### Obtain a token

To obtain a token from the newly created source:

	a3sctl auth remote-a3s \
		--source-name my-remote-a3s-source \
		--source-namespace /tutorial \
		--input-token <token>

#### Amazon STS

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

#### Google Cloud Platform token

This authentication source does not need custom source creation as it uses GCP
broadly. How to retrieve a token from GCP is beyond the scope of this document.
However, if you run a3sctl from a GCP instance, it will retrieve one for you, if
you don't pass any additional information

If you are not running the command on GCP:

	a3sctl auth gcp --access-token <token>

However, if you are running it from an GCP instance, you just need to run:

	a3sctl auth gcp

####  Azure token

This authentication source does not need custom source creation as it uses Azure
broadly. How to retrieve a token from Azure is beyond the scope of this document.
However, if you run a3sctl from a Azure instance, it will retrieve one for you, if
you don't pass any additional information

If you are not running the command on GCP:

	a3sctl auth azure --access-token <token>

However, if you are running it from an Azure instance, you just need to run:

	a3sctl auth azure

#### A3S local identity token

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

## Writing authorizations

The Authorizations allows to match a set of users (subjects) based on a claim expression and
assign them permissions. Authorizations work on white list model. Everything that
is not explicitely allowed is forbidden.

### Subject

A matching expression can be described as a basic boolean sequence like
`(org=acme && group=finance) || group=admin`. They are represented by a
2-dimensional array. As such, the expression above is written:

	[
		[ "org=admin", "group=finance" ],
		[ "group=admin" ]
	]

The first dimension represents `or` clauses and the second represents `and`
clauses.

As there are many source of authorizations and delivered claims can overlap,
potentially given way broader permissions than expected, the identity token
always contains additional claims allowing to discriminate bearer based on the
authentication source they used.

* `@source:type`: The type of source that was used to deliver the token.
* `@source:namespace`: The namespace of the source that was used.
* `@source:name`: The name of the source.

This way, you can differentiate `name=bob` based on which Bob we are aiming. A
safe subject to use in that case:

	[
		["@source:type=ldap", "@source:namespace=/my/ns", "name=bob"]
	]

This way, the authorization will only match Bob that got a token from any LDAP
authentication source that has been declared in `/my/ns`. Another Bob from
another namespace or coming from an OIDC source will not match.

### Permissions

Authorizations have then a set of permissions that describes what the matching
bearers can do. They are generic (ie they don't make assumptions about the
underlying protocol you are using) and are represented by a string of the form:

	"resource:action1,...,actionN[:id2,...idN]"

For instance, this allows bearer to walk and pet the dogs:

	"dogs:pet,walk"

This allows bearer to GET /admin:

	"/admin:get"

This allows to get and put authorizations with ID 1 or 2:

	"authorizations:get,put:1,2"

Permissions can use the `*` as resource or actions to match any. As such, the
following permission gives the bearer admin access:

	"*:*"

An authorization contains an array of permissions, granting the bearer the
union of them. If multiple authorizations match the bearer identity token, then
the union of all their permissions will be granted.

### Target namespaces

An authorization lives in a nanmespace and can target the current namespace of
some of their children. Authorizations propagate down the namespace hierarchy
starting from where it applied. It can not affect parents or sibling namespaces.

### Examples

We can create the authorizations describe above with the following command:

	a3sctl api create authorization
		--namespace /my/namespace \
		--name my-auth \
		--target-namespaces '["/my/namespace/app1"]' \
		--subject '[
			[
				"@source:type=oidc",
				"@source:namespace=/my/namespace",
				"org=admin",
				"group=finance",
			],
			[
				"@source:type=mtls",
				"@source:namespace=/my",
				"@source:name=admins",
				"group=admin",
			]
		]' \
		--permissions '["dogs:pet,walk"]'

> NOTE: If you ommit target-namespace, then the authorization applies to its own
> namespace and children.

## Check for permissions from your app

To verify a token bearer is allowed to performed some actions. The easiest way
to implement this is to add an authentication middleware in whatever HTTP
framework you are using to call a3s to verify a token and its permissions. This
middleware can call the all-in-one check endpoint `/authz`. The following
example uses curl, but you should use the HTTP communication layer currently
used in your application.

	curl -H "Content-Type: application/json" \
		-d '{
			"token": <token>,
			"resource": "/dogs"
			"action": "walk",
			"namespace: /application/namespace",
			"audience": "my-app",
		}' \
		https://127.0.0.1:44443/authz

This would return `204` if the bearer is allowed to walk the dogs in
`/application/namespace`, or `403` if either the token is invalid or the bearer
is not allowed to perform such action.

This method is the simplest but have a few drawbacks. For instance, you will
make a3s validate the token everytime, you need to make a call everytime,
and you need to transmit the bearer token at every call.

A more optimized method will be described here soon, that allows to:

* Validate token signature yourself locally
* Retrieve the entire permissions set for a given token for caching
* Validate the permissions locally
* Be notified when cached permissions needs to be invalidated.

> NOTE: This method requires the third party application to be able to
> connect to the push channel, and hence will require to be authenticated.

## Using a3sctl

a3sctl is the command line that allows to use a3s API in a user friendly manner.
It abstracts the ReST api and is self documenting. You can always get additional
help by passing the flags `--help` (or `-h`) in any command or sub command.

### Completion

a3sctl supports auto completion:

#### Bash

	. <(a3sctl completion bash)

#### Zsh

	compdef _a3sctl a3sctl
	. <(a3sctl completion zsh)

#### Fish

	. <(a3sctl completion fish)

### Configuration file

a3sctl can read the values of its flags from various places, and in that order:

* A flag directly provided
* Env variable (ie `$A3SCTL_SOURCE_NAME` for `--source-name`)
* The config file (default: `~/.config/a3sctl/default.yaml`)

You can choose the config file to user by setting the full path of the file
using the flag `--config`.

You can also pass the name of the config, without it's folder or it's extension
through the variable `A3SCTL_CONFIG_NAME`. a3sctl will scan the following
folder, in that order, to find a configuration file matching the name:

* `~/.config/a3sctl/`
* `/usr/local/etc/a3sctl/`
* `/etc/a3sctl/`

### Auto authentication

In addition to one-to-one mapping of a3sctl flags in the config file, you can
also add the key `autoauth` to automatically retrieve, cache, reuse and renew a
token using a particular authentication source. This method works for mtls and
ldap.

For instance, in `~/.config/a3sctl/default.yaml`:

	api: https://127.0.0.1:44443
	namespace: /

	autoauth:
		enable: mtls
		ldap:
			user: okenobi
			pass: '-'
			source:
				name: root
				namespace: /
		mtls:
			cert: /path/to/user-cert.pem
			key: /path/to/user-key.pem
			pass: '-'
			source:
				name: root
				namespace: /

You can decide which source to use for auto authentication by setting the
`enable` flag. Leave it empty to disable auto auth.

The token is cached in `$XDG_HOME_CACHE/a3sctl/token-<src>-<api-hash>` and will
automatically renew if it's past its half-life.

> NOTE: using `-` for secrets will automatically prompt the user for input during
> retrieval or renewing of the token.

## Dev environment

### Prerequesites

First, clone this repository and make sure you have the following installed:

* go
* mongodb
* nats-server
* tmux & tmuxinator

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
