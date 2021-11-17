# Model
model:
  rest_name: oidcsource
  resource_name: oidcsources
  entity_name: OIDCSource
  package: a3s
  group: authn
  description: An OIDC Auth source can be used to issue tokens based on existing OIDC accounts.
  get:
    description: Get a particular oidcsource object.
  update:
    description: Update a particular oidcsource object.
  delete:
    description: Delete a particular oidcsource object.
  extends:
  - '@sharded'
  - '@identifiable'

# Indexes
indexes:
- - namespace
  - name

# Attributes
attributes:
  v1:
  - name: certificateAuthority
    description: |-
      The Certificate authority to use to validate the authenticity of the OIDC
      server. If left empty, the system trust stroe will be used. In most of the
      cases, you don't need to set this.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: |-
      -----BEGIN CERTIFICATE-----
      MIIBZTCCAQugAwIBAgIRANYvXLTa16Ykvc9hQ4BBLJEwCgYIKoZIzj0EAwIwEjEQ
      MA4GA1UEAxMHQUNNRSBDQTAeFw0yMTExMDEyMzAwMTlaFw0zMTA5MTAyMzAwMTla
      MBIxEDAOBgNVBAMTB0FDTUUgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASa
      7wknroxwB1znupZ67NzTG9Kuc+tNRlbI22eTDNMKYpIexzWDOyiQ95N3GQIdmAz5
      wVu9l2V3VuKUpD9mNgkRo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUw
      AwEB/zAdBgNVHQ4EFgQURIT2kL76vMj9A3r9AUnaiHnHf4EwCgYIKoZIzj0EAwID
      SAAwRQIgS4SGaJ/B1Ul88Jal11Q5BwiY9bY2y9w+4xPNBxSyAIcCIQCSWVq+00xS
      bOmROq+EsxO4L/GzJx7MBbeJ6x142VKSBQ==
      -----END CERTIFICATE-----
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
