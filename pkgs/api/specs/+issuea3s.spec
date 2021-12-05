# Model
model:
  rest_name: issuea3s
  resource_name: issuea3s
  entity_name: IssueA3S
  package: a3s
  group: authn/issue
  description: Additional issuing information for A3S token source.
  detached: true

# Attributes
attributes:
  v1:
  - name: token
    description: The original token.
    type: string
    exposed: true
    required: true
    example_value: valid.jwt.token
