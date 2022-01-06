# Model
model:
  rest_name: httpsource
  resource_name: httpsources
  entity_name: HTTPSource
  package: a3s
  group: authn/source
  description: A source that can call a remote service to validate generic credentials.
  get:
    description: Get a particular httpsource object.
  update:
    description: Update a particular httpsource object.
  delete:
    description: Delete a particular httpsource object.
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
  - name: CA
    description: The certificate authority to use to validate the remote http server.
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

  - name: URL
    description: |-
      URL of the remote service. This URL will receive a POST containing the
      credentials information that must be validated. It must reply with 200 with a
      body containing a json array that will be used as claims for the token. Any
      other error code will be returned as a 401 error.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: https://account.acme.com/auth
    validations:
    - $url

  - name: certificate
    description: |-
      Client certificate required to call URL. A3S will refuse to send data if the
      endpoint does not support client certificate authentication.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: |-
      -----BEGIN CERTIFICATE-----
      MIIBPzCB5qADAgECAhEAwbx3c+QW24ePXyD94geytzAKBggqhkjOPQQDAjAPMQ0w
      CwYDVQQDEwR0b3RvMB4XDTE5MDIyMjIzNDA1MFoXDTI4MTIzMTIzNDA1MFowDzEN
      MAsGA1UEAxMEdG90bzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJi6CwRDeKks
      Xb3pDEslmFGR7k9Aeh5RK+XmdqKKPGb3NQWEFPGolnqOR34iVuf7KSxTuzaaVWfu
      XEa94faUQEqjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MAoG
      CCqGSM49BAMCA0gAMEUCIQD+nL9RF9EvQXHyYuJ31Lz9yWd9hsK91stnpAs890gS
      /AIgQIKjBBpiyQNZZWso5H04qke9QYMVPegiQQufFFBj32c=
      -----END CERTIFICATE-----
    validations:
    - $pem

  - name: description
    description: The description of the object.
    type: string
    exposed: true
    stored: true

  - name: key
    description: Key associated to the client certificate.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: |-
      -----BEGIN PRIVATE KEY-----
      MIIBPzCB5qADAgECAhEAwbx3c+QW24ePXyD94geytzAKBggqhkjOPQQDAjAPMQ0w
      CwYDVQQDEwR0b3RvMB4XDTE5MDIyMjIzNDA1MFoXDTI4MTIzMTIzNDA1MFowDzEN
      MAsGA1UEAxMEdG90bzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJi6CwRDeKks
      Xb3pDEslmFGR7k9Aeh5RK+XmdqKKPGb3NQWEFPGolnqOR34iVuf7KSxTuzaaVWfu
      XEa94faUQEqjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MAoG
      CCqGSM49BAMCA0gAMEUCIQD+nL9RF9EvQXHyYuJ31Lz9yWd9hsK91stnpAs890gS
      /AIgQIKjBBpiyQNZZWso5H04qke9QYMVPegiQQufFFBj32c=
      -----END PRIVATE KEY-----
    validations:
    - $pem

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
    example_value: my-http-source
