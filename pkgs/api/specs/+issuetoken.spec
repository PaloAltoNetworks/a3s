# Model
model:
  rest_name: issuetoken
  resource_name: issuetoken
  entity_name: IssueToken
  package: a3s
  group: authn
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
