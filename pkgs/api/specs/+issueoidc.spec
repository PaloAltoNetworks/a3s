# Model
model:
  rest_name: issueoidc
  resource_name: issueoidc
  entity_name: IssueOIDC
  package: a3s
  group: authn/issue
  description: Additional issuing information for the OIDC source.
  detached: true

# Attributes
attributes:
  v1:
  - name: authURL
    description: Contains the auth URL is noAuthRedirect is set to true.
    type: string
    exposed: true
    read_only: true
    omit_empty: true

  - name: code
    description: OIDC ceremony code.
    type: string
    exposed: true

  - name: noAuthRedirect
    description: |-
      If set, instruct the server to return the OIDC auth url in authURL instead of
      performing an HTTP redirection.
    type: boolean
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
