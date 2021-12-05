# Model
model:
  rest_name: issueremotea3s
  resource_name: issueremotea3s
  entity_name: IssueRemoteA3S
  package: a3s
  group: authn/issue
  description: Additional issuing information for a remote A3S token source.
  detached: true

# Attributes
attributes:
  v1:
  - name: token
    description: The remote a3s token.
    type: string
    exposed: true
    required: true
    example_value: valid.jwt.token
