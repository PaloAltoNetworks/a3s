# Model
model:
  rest_name: authz
  resource_name: authz
  entity_name: Authz
  package: a3s
  group: authz/check
  description: API to verify permissions.

# Attributes
attributes:
  v1:
  - name: ID
    description: The optional ID of the object to check permission for.
    type: string
    exposed: true

  - name: IP
    description: IP of the client.
    type: string
    exposed: true

  - name: action
    description: The action to check permission for.
    type: string
    exposed: true
    required: true
    example_value: delete

  - name: audience
    description: Audience that should be checked for.
    type: string
    exposed: true

  - name: namespace
    description: The namespace where to check permission from.
    type: string
    exposed: true
    required: true
    example_value: /acme

  - name: resource
    description: The resource to check permission for.
    type: string
    exposed: true
    required: true
    example_value: cats

  - name: token
    description: The token to check.
    type: string
    exposed: true
    subtype: string
    required: true
    example_value: aaa.valid.jwt
