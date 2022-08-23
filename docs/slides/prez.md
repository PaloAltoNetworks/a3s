---
author: Antoine Mercadal
---

**pre-flight**

- a3s running
- testapp/server.py running and initialized
- ldap server running
- clean up import

```bash
a3sctl api delete ldapsource ldap -n /testapp
a3sctl api delete authorization ldap -n /testapp
```

---
# A3S

Authentication As A Service

---
## General Concepts

* Provides Authentication using **JWT**
* Provides Authorization using generic **ABAC** system
* Policies are based on JWT **identity claims**
* Everything is organized in a **tree** data structure
* **Federated** and **decentralized**
* Extracted and modernized from Aporeto Core
* Designed to be **fast** and **scalable**

---
## Namespaces

* Namespaces are the most basic construct
* Every resources live in a namespace (including namespaces)
* Policy applies to a namespace and propagates down

Example namespaces hierarchy:

```
~~~graph-easy --as boxart
[ / ] -> { start: front, 0; } [ /coca ]
[ /coca ] {rows: 5} -> { start: front, 0; } [ /coca/qa ], [ /coca/dev ], [ /coca/prod ]
[ / ] -> { start: front, 0; } [ /pepsi ]
[ /pepsi ] {rows: 5} -> { start: front, 0; } [ /pepsi/qa ], [ /pepsi/dev ], [ /pepsi/prod ]
~~~
```

---
## Authentication

* Users declare **external sources of authentications** in namespaces
* Clients can use any source they have credentials for and get a **normalized token**
* Identity is **derived** from the source used as well as data extracted from the source

> The JWT are never used to declare permissions. They are like an ID card, stating
who the bearer is. Not what they can do.

> The identity of a bearer does not change over time. Only their permissions do.

---
### Sources

A3S supports the following authentication sources:

- MTLS (X.509)
- OIDC
- LDAP
- Azure, GCP, AWS identity tokens
- SAML
- Generic HTTP call
- Other A3S instances (federation)

---
#### Identity token from a X.509 Certificate

```json
{
  "aud": [
    "https://api.system1.com"
  ],
  "exp": 1651268495,
  "iat": 1651182095,
  "identity": [
    "@issuer=https://api.system1.com",
    "@source:name=root",
    "@source:namespace=/",
    "@source:type=mtls",
    "akid=4484F690BEFABCC8FD037AFD0149DA8879C77F81",
    "commonname=Jean-Michel",
    "fingerprint=52678D7B924E7A0D40C9FE35FE049909731014F1D86E09F3898A84D884087C6B",
    "issuerchain=93DC074396E4066BA01C126B6136FFD530197AA2FC87C42FFB28FBE0F524DF8C",
    "serialnumber=219959457279438724775594138274989969558"
  ],
  "iss": "https://api.system1.com",
  "jti": "e9c1eba4-227a-4690-98be-a2c7b0e09294"
}
```

---
#### Identity token from OIDC

```json
{
  "aud": [
    "https://api.system1.com"
  ],
  "exp": 1651187692,
  "iat": 1651184092,
  "identity": [
    "@issuer=https://api.system1.com",
    "@source:name=gcp",
    "@source:namespace=/",
    "@source:type=oidc",
    "email=amercadal@paloaltonetworks.com",
    "email_verified=true",
    "hd=paloaltonetworks.com",
    "iss=https://accounts.google.com"
  ],
  "iss": "https://api.system1.com",
  "jti": "72d1fe55-d146-4cd8-8af1-123cb710f835"
}
```

---
#### Identity token from LDAP

```json
{
  "aud": [
    "https://api.system1.com"
  ],
  "exp": 1651269229,
  "iat": 1651182829,
  "identity": [
    "@issuer=https://api.system1.com",
    "@source:name=a3s-dev-ldap",
    "@source:namespace=/",
    "@source:type=ldap",
    "cn=User1",
    "cn=okenobi",
    "dc=io",
    "dc=universe",
    "dn=cn=okenobi,ou=users,dc=universe,dc=io",
    "gidNumber=1000",
    "homeDirectory=/home/okenobi",
    "ou=users",
    "sn=Bar1",
    "uid=okenobi",
    "uidNumber=1000"
  ],
  "iss": "https://api.system1.com",
  "jti": "899aa32d-8340-4712-968c-c56231f30d75"
}
```

---
### Identity Claims

* A3S always sets the following claims:

```text
@source:name=x
@source:namespace=/namespace
@source:type=y
@issuer=z
```

* All other claims are derived from the source

```text
name=john
team=red
team=blue
email=john@domain.com
```

> A3S supports **claim cloaking**, letting the user in control of the claims they
want to share.

---
## Authorization

* Authorization policies **discriminate** users based on **identity claims**
* They use a **claim selector expression** to define on who they apply
* They define a set of **resources and actions** bearers matching the claims selector can do
* They are defined in a namespace, apply to that namespace or some below and propagate down.
* They can also understand tokens issued from other instances of A3S (federation)

---
### Selectors

A claim selector expression is a two dimensional array of ORs of ANDs.

```json
[
  [
    "@source:type=oidc",
    "@source:namespace=/tenant/users"
  ],
  [
    "@source:type=mtls",
    "@source:namespace=/tenant/admins",
    "commonname=Antoine Mercadal"
  ]
]
```

This matches a token with both `@source:type=oidc` and
`@source:namespace=/tenant/users` **or** a token with all three claims
`@source:type=mtls`, `@source:namespace=/tenant/admins` and `commonname=Antoine
Mercadal`

---
### Permissions

Permissions are simple string expressions defining some actions on a particular
resource, optionally with particular identifiers.

Generically:

```text
<resource>:<action-1>,...,<action-N>[:<id-1>,...,<id-N>]
```

Specifically:

```text
dog:walk,pet
cat:feed:kitty
/admin:get,put,post,delete
/resources/object:get:42,43
feature-x:*
*:get
```

---
### Example

```json
{
  "name": "pets-walkers",
  "description": "allow pets walking",
  "namespace": "/tenant",
  "targetNamespaces": [
    "/tenant/animals"
  ],
  "subject": [
    [
      "@source:type=oidc",
      "@source:name=default",
      "@source:namespace=/tenant/users"
    ]
  ],
  "permissions": [
    "dog:walk",
    "cat:walk"
  ],
  "subnets": null,
  "trustedIssuers": [
    "https://api.system1.com"
  ]
}
```

---
## Demo

The system is configured like the following:

```
~~~graph-easy --as boxart
[ client ] - \l authn -> [ A3S ]
[ client ] - access ->  [ testapp ]
[ testapp ] - authz ->  [ A3S ]
~~~
```

---
### TestApp

The TestApp is a simple Python Flask webapp that provides 3 endoints:

* `/`
* `/secret`
* `/topsecret`

This first one is public, while the two other requires permissions:

```python
~~~./extract.awk routes:start routes:end ../../examples/python/testapp/server.py

~~~
```

---
#### Naive authenticator

Naive implementation of an authenticator/authorizer leveraging the `/authz`
endpoint from A3S:

```python
~~~./extract.awk authenticator:start authenticator:end ../../examples/python/testapp/server.py

~~~
```

---
#### TestApp Policies

Policy imported by the TestApp administrator into A3S:

```yaml
~~~cat ../../examples/python/testapp/import.gotmpl

~~~
```

---
#### Accessing without a token

```bash
curl -k https://127.0.0.1:5000
curl -k https://127.0.0.1:5000/secret
curl -k https://127.0.0.1:5000/topsecret
```

---
#### Accessing as John

John gets a token using the configured X.509 authentication source:

```bash
token=$(
  a3sctl auth mtls \
    --api https://127.0.0.1:44443 \
    --api-skip-verify \
    --audience testapp \
    --source-namespace /testapp \
    --cert ../../examples/python/testapp/certs/john-cert.pem \
    --key ../../examples/python/testapp/certs/john-key.pem
)

curl -k -u "Bearer:$token" https://127.0.0.1:5000
curl -k -u "Bearer:$token" https://127.0.0.1:5000/secret
curl -k -u "Bearer:$token" https://127.0.0.1:5000/topsecret
```

---
#### Accessing as Michael

Michael gets a token using the configured X.509 authentication source:

```bash
token=$(
  a3sctl auth mtls \
    --api https://127.0.0.1:44443 \
    --api-skip-verify \
    --audience testapp \
    --source-namespace /testapp \
    --cert ../../examples/python/testapp/certs/michael-cert.pem \
    --key ../../examples/python/testapp/certs/michael-key.pem
)

curl -k -u "Bearer:$token" https://127.0.0.1:5000
curl -k -u "Bearer:$token" https://127.0.0.1:5000/secret
curl -k -u "Bearer:$token" https://127.0.0.1:5000/topsecret
```

---
#### Accessing the app using a browser

We first add a new LDAP source and Authorization to the /testapp namespace to
allow people login from a LDAP server:

```bash
A3SCTL_API=https://124.0.0.1:44443
A3SCTL_API_SKIP_VERIFY=1

cat <<EOF | a3sctl import -n /testapp - 2>&1
label: ldap
LDAPSources:
- name: ldap
  address: 127.0.0.1:11389
  baseDN: dc=universe,dc=io
  bindDN: cn=admin,dc=universe,dc=io
  bindPassword: "password"
  securityProtocol: None
Authorizations:
- name: ldap
  subject:
  - - "@source:type=ldap"
    - "@source:name=ldap"
    - "@source:namespace=/testapp"
    - "cn=okenobi"
  permissions:
  - "/topsecret:GET"
EOF
```

```bash
flatpak run org.mozilla.firefox --private-window https://localhost:5000?rlogin
```

---
## Additional Features

* IRL authentication/authorization based on QR codes
* Claims modifiers hooks
* Restricted tokens:
	* permissions
	* origin subnets
	* target namespaces
* Helm like templates in import files
* Audit trails

---
## Thanks!

Questions? Opinions? Insults?
