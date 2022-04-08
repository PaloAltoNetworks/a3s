# Model
model:
  rest_name: authorization
  resource_name: authorizations
  entity_name: Authorization
  package: a3s
  group: authz
  description: TODO.
  get:
    description: Retrieves the authorization with the given ID.
  update:
    description: Updates the authorization with the given ID.
  delete:
    description: Deletes the authorization with the given ID.
    global_parameters:
    - $queryable
  extends:
  - '@sharded'
  - '@identifiable'
  - '@importable'

# Indexes
indexes:
- - namespace
  - flattenedSubject
  - disabled
- - namespace
  - flattenedSubject
  - propagate
- - namespace
  - trustedIssuers

# Attributes
attributes:
  v1:
  - name: description
    description: Description of the Authorization.
    type: string
    exposed: true
    stored: true

  - name: disabled
    description: Set the authorization to be disabled.
    type: boolean
    exposed: true
    stored: true

  - name: flattenedSubject
    description: This is a set of all subject tags for matching in the DB.
    type: list
    subtype: string
    stored: true

  - name: hidden
    description: Hides the policies in children namespaces.
    type: boolean
    exposed: true
    stored: true

  - name: name
    description: The name of the Authorization.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: my authorization

  - name: permissions
    description: A list of permissions.
    type: list
    exposed: true
    subtype: string
    stored: true
    required: true
    example_value:
    - '@auth:role=namespace.administrator'
    - namespace,get,post,put
    - authorization,get:1234567890

  - name: propagate
    description: Propagates the api authorization to all of its children. This is always true.
    type: boolean
    stored: true
    default_value: true
    getter: true
    setter: true

  - name: subject
    description: A tag expression that identifies the authorized user(s).
    type: external
    exposed: true
    subtype: '[][]string'
    stored: true
    orderable: true
    validations:
    - $tags_expression
    - $authorization_subject

  - name: subnets
    description: |-
      If set, the API authorization will only be valid if the request comes from one
      the declared subnets.
    type: list
    exposed: true
    subtype: string
    stored: true
    validations:
    - $cidr_list_optional

  - name: targetNamespaces
    description: |-
      Defines the namespace or namespaces in which the permission for subject should
      apply. If empty, the object's namespace will be used.
    type: list
    exposed: true
    subtype: string
    stored: true
    example_value: /my/namespace

  - name: trustedIssuers
    description: List of issuers to consider before using the policy for a given set of claims.
    type: list
    exposed: true
    subtype: string
    stored: true
