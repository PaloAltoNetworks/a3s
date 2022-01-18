# Model
model:
  rest_name: import
  resource_name: import
  entity_name: Import
  package: a3s
  group: core/import
  description: Import multiple resource at once.

# Attributes
attributes:
  v1:
  - name: A3SSources
    description: A3S sources to import.
    type: refList
    exposed: true
    subtype: a3ssource
    omit_empty: true

  - name: HTTPSources
    description: HTTP sources to import.
    type: refList
    exposed: true
    subtype: httpsource
    omit_empty: true

  - name: LDAPSources
    description: LDAP sources to import.
    type: refList
    exposed: true
    subtype: ldapsource
    omit_empty: true

  - name: MTLSSources
    description: MTLS sources to import.
    type: refList
    exposed: true
    subtype: mtlssource
    omit_empty: true

  - name: OIDCSources
    description: OIDC sources to import.
    type: refList
    exposed: true
    subtype: oidcsource
    omit_empty: true

  - name: authorizations
    description: Authorizations to import.
    type: refList
    exposed: true
    subtype: authorization
    omit_empty: true

  - name: label
    description: |-
      Import label that will be used to identify all the resources imported by this
      resource.
    type: string
    exposed: true
    required: true
    example_value: my-super-import

  - name: mode
    description: Import mode. If set to Remove, the previously imported data will be removed.
    type: enum
    exposed: true
    allowed_choices:
    - Import
    - Remove
    default_value: Import
