# a3s dev ldap

This folder contains a ready to use pre-configured ldap server and an import
file to declare it as a source in A3S.

> Note: This entire setup is insecure. This is intended for development
> purposes.

## Start the container

	docker compose up

The server listens on port 11389, with no TLS.

* admin username: `cn=admin,dc=universe,dc=io`
* admin password: `password`

It contains 2 additional pre-populated users:

- username: `okenobi` password: `pass`
- username: `dvader` password: `pass`

## Import the ldap source

To use the ldap server as an A3S ldap source, import its declaration with the
following command:

	a3sctl api create import \
		--input-file ./a3s-ldapsource.yaml \
		--namespace /

> Note: There will be a helper subcommand to handle imports soon. The version
> above uses the raw API. Check apoctl --help.

The source name is `a3s-dev-ldap`.
