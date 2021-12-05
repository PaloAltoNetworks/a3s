# Model
model:
  rest_name: issueaws
  resource_name: issueaws
  entity_name: IssueAWS
  package: a3s
  group: authn/issue
  description: Additional issuing information for AWS STS token source.
  detached: true

# Attributes
attributes:
  v1:
  - name: ID
    description: The ID of the AWS STS token.
    type: string
    exposed: true
    required: true
    example_value: xxxxx

  - name: secret
    description: The secret associated to the AWS STS token.
    type: string
    exposed: true
    required: true
    example_value: yyyyy

  - name: token
    description: The original token.
    type: string
    exposed: true
    required: true
    example_value: valid.jwt.token
