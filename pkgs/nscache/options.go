package nscache

type config struct {
	notificationName string
}

func newConfig() config {
	return config{
		notificationName: NotificationNamespaceChanges,
	}
}

// Option represents a parametric option for the nscache.
type Option func(*config)

// OptionNotificationName allows to change the notification name
// used to listen to namespace changes from the pubsub.
// This defaults to OptionNotificationName,
func OptionNotificationName(name string) Option {
	return func(c *config) {
		c.notificationName = name
	}
}
