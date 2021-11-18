package oidcceremony

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"golang.org/x/oauth2"
)

const oidcCacheCollection = "oidccache"

// CacheItem represents a cache OIDC request info.
type CacheItem struct {
	State            string        `bson:"state"`
	ClientID         string        `bson:"clientid"`
	CA               string        `bson:"ca"`
	OAuth2Config     oauth2.Config `bson:"oauth2config"`
	ProviderEndpoint string        `bson:"providerEndpoint"`
	Time             time.Time     `bson:"time"`
}

// Set sets the given OIDCRequestItem in redis.
func Set(m manipulate.Manipulator, item *CacheItem) error {

	item.Time = time.Now()

	db, disco, err := manipmongo.GetDatabase(m)
	if err != nil {
		return err
	}
	defer disco()

	return db.C(oidcCacheCollection).Insert(item)
}

// Get gets the items with the given state.
// If none is found, it will return nil.
func Get(m manipulate.Manipulator, state string) (*CacheItem, error) {

	db, disco, err := manipmongo.GetDatabase(m)
	if err != nil {
		return nil, err
	}
	defer disco()

	item := &CacheItem{}
	if err := db.C(oidcCacheCollection).Find(bson.M{"state": state}).One(item); err != nil {
		return nil, err
	}
	return item, nil
}

// Delete deletes the items with the given state.
func Delete(m manipulate.Manipulator, state string) error {

	db, disco, err := manipmongo.GetDatabase(m)
	if err != nil {
		return err
	}
	defer disco()

	return db.C(oidcCacheCollection).Remove(bson.M{"state": state})
}
