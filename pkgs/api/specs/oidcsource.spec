# Model
model:
  rest_name: oidcsource
  resource_name: oidcsources
  entity_name: OIDCSource
  package: a3s
  group: authn/source
  description: An OIDC Auth source can be used to issue tokens based on existing OIDC
    accounts.
  get:
    description: Get a particular oidcsource object.
  update:
    description: Update a particular oidcsource object.
  delete:
    description: Delete a particular oidcsource object.
  extends:
  - '@sharded'
  - '@identifiable'
  - '@importable'

# Indexes
indexes:
- - namespace
  - name

# Attributes
attributes:
  v1:
  - name: CA
    description: |-
      The Certificate authority to use to validate the authenticity of the OIDC
      server. If left empty, the system trust stroe will be used. In most of the
      cases, you don't need to set this.
    type: string
    exposed: true
    stored: true
    validations:
    - $pem

  - name: clientID
    description: Unique client ID.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: 12345677890.apps.googleusercontent.com

  - name: clientSecret
    description: Client secret associated with the client ID.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: Ytgbfjtj4652jHDFGls99jF
    encrypted: true

  - name: description
    description: The description of the object.
    type: string
    exposed: true
    stored: true

  - name: endpoint
    description: |-
      OIDC [discovery
      endpoint](https://openid.net/specs/openid-connect-discovery-1_0.html#IssuerDiscovery).
    type: string
    exposed: true
    stored: true
    required: true
    example_value: https://accounts.google.com

  - name: modifier
    description: |-
      Contains optional information about a remote service that can be used to modify
      the claims that are about to be delivered using this authentication source.
    type: ref
    exposed: true
    subtype: identitymodifier
    stored: true
    omit_empty: true
    extensions:
      noInit: true
      refMode: pointer

  - name: name
    description: The name of the source.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: myoidc

  - name: scopes
    description: List of scopes to allow.
    type: list
    exposed: true
    subtype: string
    stored: true
    example_value:
    - email
    - profile
