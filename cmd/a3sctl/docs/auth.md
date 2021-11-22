## auth

The auth subcommand allows to retrieve an identity token from one of the various
support authentication sources:

* `a3s`: exchange an existing a3s token for a new one.
* `aws`: use AWS STS.
* `azure`: use Azure identity token.
* `check`: check the content of a token.
* `gcp`: use Google Cloud Platfrom identity token.
* `ldap`: use a configured LDAP authentication source.
* `mtls`: use a client certificate.
* `oidc`: use a configured OIDC authentication source.

### Validity

You can ask for a custom token validity with the flag `--validity`. The server
may decide to cap the requested to an admin configured maximum life time.

## Restrictions

You can restrict the delivered token by using the following flags:

* `--restrict-namespace`: the delivered token will only be valid for the
    restricted namespace.
* `--restrict-network`: the delivered token will only be valid if used from one
    of the provided network. This network must be visible as an origin from the
    a3s server. If you use NAT, you need to use the NAT'ed address.
* `--restrict-permissions`: Whatever is the set of permissions allowed per
    authorization policies, restrict to only the provided permissions.
    Restricted permissions must be contained into the actualy granted
    permissions set.

### Cloaking

You may want to limit the number of identity claims that will be embeded into
the delivered token. Cloaking allows to pass claim prefixes that must match in
order for it to be added in the token. For intance, to ask for a token only
containing your age, you can use `--cloak age=`.

### Audience

You can ask a3s to deliver a token for various audience. The flag `--audience`
can be used multiple time to set the desired audience. For example:

    --audience myapp1 --audience myapp2

### Auto authentication

You can configure a3sctl to automatically use a source of your choice to
retrieve a token and renew it when needed. To do so, you need to add a section
named `autoauth` in the a3sctl configuration file.

For now, a3sctl only support stable auth source like LDAP or MTLS.

For example:

    autoauth:
      enable: mtls
      mtls:
        cert: /path/to/user-cert.pem
        key: /path/to/user-key.pem
        pass: '-'
        source:
          name: root
          namespace: /
