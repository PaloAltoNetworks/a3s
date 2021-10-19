# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: root
  group: core
  description: root object.
  root: true

# Relations
relations:
- rest_name: issue
  create:
    description: Ask to issue a new authentication token.
    parameters:
      entries:
      - name: asCookie
        description: If set to true, the token will be delivered in a secure cookie, and not in the response body.
        type: boolean

- rest_name: namespace
  get:
    description: Retrieves the list of namespaces.
  create:
    description: Creates a new namespace.
