# Model
model:
  rest_name: namespace
  resource_name: namespaces
  entity_name: Namespace
  package: a3s
  group: policy
  description: |-
    A namespace is grouping object. Every object is part of a namespace, and every
    request is made against a namespace. Namespaces form a tree hierarchy.
  get:
    description: Get a particular namespace object.
  update:
    description: Update a particular namespace object.
  delete:
    description: Delete a particular namespace object.
  extends:
  - '@sharded'
  - '@identifiable'

# Indexes
indexes:
- - namespace
  - name
- - name

# Attributes
attributes:
  v1:
  - name: description
    description: The description of the object.
    type: string
    exposed: true
    stored: true

  - name: name
    description: |-
      The name of the namespace. When you create a namespace, only put its bare name,
      not its full path.
    type: string
    exposed: true
    stored: true
    required: true
    creation_only: true
    allowed_chars: ^[a-zA-Z0-9_/]+$
    allowed_chars_message: must only contain alpha numerical characters, '-' or '_'
    example_value: mycompany
    getter: true
    setter: true
