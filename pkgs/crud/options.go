package crud

import (
	"fmt"

	"go.aporeto.io/elemental"
)

// A ErrPreWriteHook is the kind of error returned when execution
// of a pre write hook fails.
type ErrPreWriteHook struct {
	Err error
}

func (e ErrPreWriteHook) Error() string {
	return fmt.Sprintf("unable to run pre-write hook: %s", e.Err)
}

func (e ErrPreWriteHook) Unwrap() error {
	return e.Err
}

// PreWriteHook is the type of function you can pass as a pre write hook.
type PreWriteHook func(obj elemental.Identifiable, orig elemental.Identifiable) error

// PostWriteHook is the type of function yiou can pass as a post write hook.
type PostWriteHook func(obj elemental.Identifiable)

type cfg struct {
	preHook  PreWriteHook
	postHook PostWriteHook
}

// An Option defines optional behaviors for the crud functions.
type Option func(*cfg)

// OptionPreWriteHook defines the function to call just before the final
// write operation runs. It will be given the object that is about to be
// written and eventually the original object that got pulled from the database.
// If this function returns an error, the operation stops and the error will
// be returned wrapped into a PreWriteHookError.
func OptionPreWriteHook(hook PreWriteHook) Option {
	return func(c *cfg) {
		c.preHook = hook
	}
}

// OptionPostWriteHook defines the function to call just after the final
// write operation runs. It will be given the object that was written.
func OptionPostWriteHook(hook PostWriteHook) Option {
	return func(c *cfg) {
		c.postHook = hook
	}
}
