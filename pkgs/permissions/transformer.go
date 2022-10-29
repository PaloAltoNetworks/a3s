package permissions

// A Transformer is an interface that can manipulate a permissions map.
type Transformer interface {
	Transform(PermissionMap) PermissionMap
}
