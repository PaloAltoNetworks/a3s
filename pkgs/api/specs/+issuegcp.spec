# Model
model:
  rest_name: issuegcp
  resource_name: issuegcp
  entity_name: IssueGCP
  package: a3s
  group: authn
  description: Additional issuing information for GCP identity token source.
  detached: true

# Attributes
attributes:
  v1:
  - name: audience
    description: The required audience.
    type: string
    exposed: true

  - name: token
    description: The original token.
    type: string
    exposed: true
    required: true
    example_value: valid.jwt.token
