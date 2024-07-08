# Model
model:
  rest_name: namespacedeletionrecord
  resource_name: namespacedeletionrecords
  entity_name: NamespaceDeletionRecord
  package: a3s
  group: core
  description: |-
    A namespace deletion record holds the namespace that was deleted and the date it
    was deleted.
  extends:
  - '@sharded'
  - '@identifiable'

# Attributes
attributes:
  v1:
  - name: deleteTime
    description: Deletion date of the object.
    type: time
    exposed: true
    stored: true

  - name: namespace
    description: Namespace that got deleted.
    type: string
    exposed: true
    stored: true
