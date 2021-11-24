# Model
model:
  rest_name: identitymodifier
  resource_name: identitymodifier
  entity_name: IdentityModifier
  package: a3s
  group: authn
  description: |-
    Information about a remote endpoint to call to eventually modify the identity
    claims about to be issued when using the parent source.
  detached: true

# Attributes
attributes:
  v1:
  - name: URL
    description: |-
      URL of the remote service. This URL will receive a call containing the
      claims that are about to be delivered. It must reply with 204 if it does not
      wish to modify the claims, or 200 alongside a body containing the modified
      claims.
    type: string
    exposed: true
    required: true
    example_value: https://modifier.acme.com/modify

  - name: certificate
    description: |-
      Client certificate required to call URL. A3S will refuse to send data if the
      endpoint does not support client certificate authentication.
    type: string
    exposed: true
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

  - name: certificateAuthority
    description: CA to use to validate the entity serving the URL.
    type: string
    exposed: true
    stored: true
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
    omit_empty: true

  - name: key
    description: Key associated to the client certificate.
    type: string
    exposed: true
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

  - name: method
    description: |-
      The HTTP method to use to call the endpoint. For POST/PUT/PATCH the remote
      server will receive the claims as a JSON encoded array in the body. For a GET, the claims will be passed as a query parameter named `claim`.
    type: enum
    exposed: true
    required: true
    allowed_choices:
    - GET
    - POST
    - PUT
    - PATCH
    default_value: POST
