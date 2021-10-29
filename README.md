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
