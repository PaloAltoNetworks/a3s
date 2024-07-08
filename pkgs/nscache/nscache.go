package nscache

import (
	"context"
	"time"

	"github.com/karlseguin/ccache/v2"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/bahamut"
)

// Constants for notification topics.
const (
	NotificationNamespaceChanges = "notifications.changes.namespace"
)

// A NamespacedCache is used to cache namespaced information.
// The cache will invalidate all items when their namespace is
// deleted or updated.
type NamespacedCache struct {
	pubsub           bahamut.PubSubClient
	cache            *ccache.Cache
	notificationName string
}

// New returns a new namespace cache.
func New(pubsub bahamut.PubSubClient, maxSize int64, options ...Option) *NamespacedCache {

	cfg := newConfig()
	for _, o := range options {
		o(&cfg)
	}

	return &NamespacedCache{
		pubsub:           pubsub,
		cache:            ccache.New(ccache.Configure().MaxSize(maxSize)),
		notificationName: cfg.notificationName,
	}
}

// Set sets a new namespaced key with the given value, with given expiration.
// namespace must be set. key is optional. It can be empty if you wish to only associate
// one value to one namespace.
func (c *NamespacedCache) Set(namespace string, key string, value any, duration time.Duration) {

	c.cache.Set(namespace+":"+key, value, duration)
}

// Get returns the cached item for the provided namespaced key.
func (c *NamespacedCache) Get(namespace string, key string) *ccache.Item {

	return c.cache.Get(namespace + ":" + key)
}

// Delete attempts to delete an item from the cache using the given namespace and key.
func (c *NamespacedCache) Delete(namespace string, key string) bool {

	return c.cache.Delete(namespace + ":" + key)
}

// Start starts listening to notifications for automatic invalidation
func (c *NamespacedCache) Start(ctx context.Context) {

	notification.Subscribe(
		ctx,
		c.pubsub,
		c.notificationName,
		func(msg *notification.Message) {
			c.cleanupCacheForNamespace(msg.Data.(string))
		},
	)
}

func (c *NamespacedCache) cleanupCacheForNamespace(ns string) {

	suffix := "/"
	if ns == "/" {
		suffix = ""
	}

	c.cache.DeletePrefix(ns + ":")
	c.cache.DeletePrefix(ns + suffix)
}
