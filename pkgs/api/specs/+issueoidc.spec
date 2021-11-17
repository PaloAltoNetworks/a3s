# Model
model:
  rest_name: issueoidc
  resource_name: issueoidc
  entity_name: IssueOIDC
  package: a3s
  group: authn
  description: Additional issuing information for the OIDC source.
  detached: true

# Attributes
attributes:
  v1:
  - name: code
    description: OIDC ceremony code.
    type: string
    exposed: true

  - name: redirectErrorURL
    description: OIDC redirect url in case of error.
    type: string
    exposed: true

  - name: redirectURL
    description: OIDC redirect url.
    type: string
    exposed: true

  - name: state
    description: OIDC ceremony state.
    type: string
    exposed: true
