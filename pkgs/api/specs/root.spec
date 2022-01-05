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
- rest_name: a3ssource
  get:
    description: Retrieves the list of a3ssources.
    global_parameters:
    - $queryable
  create:
    description: Creates a new a3ssource.

- rest_name: authorization
  get:
    description: Retrieves the list of authorization.
    global_parameters:
    - $queryable
  create:
    description: Creates a new authorization.

- rest_name: authz
  create:
    description: Sends a authz request.

- rest_name: httpsource

  get:
    description: Retrieves the list of httpsources.
    global_parameters:
    - $queryable
  create:
    description: Creates a new httpsource.

- rest_name: issue
  create:
    description: Ask to issue a new authentication token.

- rest_name: ldapsource
  get:
    description: Retrieves the list of ldapsources.
    global_parameters:
    - $queryable
  create:
    description: Creates a new ldapsource.

- rest_name: mtlssource
  get:
    description: Retrieves the list of mtlssources.
    global_parameters:
    - $queryable
  create:
    description: Creates a new mtlssource.

- rest_name: namespace
  get:
    description: Retrieves the list of namespaces.
    global_parameters:
    - $queryable
  create:
    description: Creates a new namespace.

- rest_name: oidcsource
  get:
    description: Retrieves the list of oidcsources.
    global_parameters:
    - $queryable
  create:
    description: Creates a new oidcsource.

- rest_name: permissions
  create:
    description: Sends a permissions request.
