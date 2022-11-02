package timeable

import "time"

// A Timeable is an entity that holds a create and update time.
type Timeable interface {
	GetCreateTime() time.Time
	SetCreateTime(time.Time)
	GetUpdateTime() time.Time
	SetUpdateTime(time.Time)
}
