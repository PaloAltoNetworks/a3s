# Model
model:
  rest_name: issue
  resource_name: issue
  entity_name: Issue
  package: authn
  group: authn/issue
  description: Issues a new a normalized token using various authentication sources.
  validations:
  - $issue

# Attributes
attributes:
  v1:
  - name: audience
    description: Requested audience for the delivered token.
    type: list
    exposed: true
    subtype: string
    example_value:
    - https://myfirstapp
    - https://mysecondapp
    omit_empty: true

  - name: cloak
    description: |-
      Sets a list of identity claim prefix to allow in the final token. This can be
      used to hide some information when asking for a token as not all systems need to
      know all of the claims.
    type: list
    exposed: true
    subtype: string
    example_value:
    - org=
    - age=
    omit_empty: true

  - name: cookie
    description: If set, return the token as a secure cookie.
    type: boolean
    exposed: true
    omit_empty: true

  - name: cookieDomain
    description: If set, use the provided domain for the delivered cookie.
    type: string
    exposed: true
    omit_empty: true

  - name: inputA3S
    description: Contains additional information for an A3S token source.
    type: ref
    exposed: true
    subtype: issuea3s
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputAWS
    description: Contains additional information for an AWS STS token source.
    type: ref
    exposed: true
    subtype: issueaws
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputAzure
    description: Contains additional information for an Azure token source.
    type: ref
    exposed: true
    subtype: issueazure
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputGCP
    description: Contains additional information for an GCP token source.
    type: ref
    exposed: true
    subtype: issuegcp
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputHTTP
    description: Contains additional information for an HTTP source.
    type: ref
    exposed: true
    subtype: issuehttp
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputLDAP
    description: Contains additional information for an LDAP source.
    type: ref
    exposed: true
    subtype: issueldap
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputOIDC
    description: Contains additional information for an OIDC source.
    type: ref
    exposed: true
    subtype: issueoidc
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: inputRemoteA3S
    description: Contains additional information for a remote A3S token source.
    type: ref
    exposed: true
    subtype: issueremotea3s
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: opaque
    description: Opaque data that will be included in the issued token.
    type: external
    exposed: true
    subtype: map[string]string
    omit_empty: true

  - name: restrictedNamespace
    description: |-
      Restricts the namespace where the token can be used.

      For instance, if you have have access to `/namespace` and below, you can
      tell the policy engine that it should restrict further more to
      `/namespace/child`.

      Restricting to a namespace you don't have initially access according to the
      policy engine has no effect and may end up making the token unusable.
    type: string
    exposed: true
    example_value: /namespace
    omit_empty: true

  - name: restrictedNetworks
    description: |-
      Restricts the networks from where the token can be used. This will reduce the
      existing set of authorized networks that normally apply to the token according
      to the policy engine.

      For instance, If you have authorized access from `0.0.0.0/0` (by default) or
      from
      `10.0.0.0/8`, you can ask for a token that will only be valid if used from
      `10.1.0.0/16`.

      Restricting to a network that is not initially authorized by the policy
      engine has no effect and may end up making the token unusable.
    type: list
    exposed: true
    subtype: string
    example_value:
    - 10.0.0.0/8
    - 127.0.0.1/32
    omit_empty: true
    validations:
    - $cidr_list_optional

  - name: restrictedPermissions
    description: |-
      Restricts the permissions of token. This will reduce the existing permissions
      that normally apply to the token according to the policy engine.

      For instance, if you have administrative role, you can ask for a token that will
      tell the policy engine to reduce the permission it would have granted to what is
      given defined in the token.

      Restricting to some permissions you don't initially have according to the policy
      engine has no effect and may end up making the token unusable.
    type: list
    exposed: true
    subtype: string
    example_value:
    - dogs,post
    omit_empty: true

  - name: sourceName
    description: The name of the source to use.
    type: string
    exposed: true
    example_value: /my/ns
    omit_empty: true

  - name: sourceNamespace
    description: The namespace of the source to use.
    type: string
    exposed: true
    example_value: /my/ns
    omit_empty: true

  - name: sourceType
    description: |-
      The authentication source. This will define how to verify
      credentials from internal or external source of authentication.
    type: enum
    exposed: true
    required: true
    allowed_choices:
    - A3S
    - AWS
    - Azure
    - GCP
    - HTTP
    - LDAP
    - MTLS
    - OIDC
    - RemoteA3S
    - SAML
    example_value: OIDC

  - name: token
    description: Issued token.
    type: string
    exposed: true
    read_only: true
    autogenerated: true
    omit_empty: true

  - name: validity
    description: |-
      Configures the maximum length of validity for a token, using
      [Golang duration syntax](https://golang.org/pkg/time/#example_Duration). If it
      is bigger than the configured max validity, it will be capped. Default: `24h`.
    type: string
    exposed: true
    default_value: 24h
    omit_empty: true
    validations:
    - $duration
