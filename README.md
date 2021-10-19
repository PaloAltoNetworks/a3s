# a3s

> NOTE: this is a work in progress and this software is not usable yet

a3s (stands for Auth As A Service) is an authentication and ABAC authorization
server.

It allows to normalize various sources of authentication like OIDC,
AWS/Azure/GCP Identity tokens into a generic authentication token that contains
identity claims. These claims can be used in some authorization policies to give
a particular subset of users various permissions.

These authorization policies are used to match a set of users based on a logical
claim expression (like `group=red and color=blue or group=admin`) and apply to a
namespace.

A namespace is an organizational node that is part of hierarchical tree.

Basically, an authorization policy allows a subset of users based on the claims
retrieved from an authentication source to perform actions in a particular
namespace and all of its children.
