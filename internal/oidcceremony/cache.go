package oidcceremony

import (
	"context"
	"time"

	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

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

	db := manipmongo.GetDatabase(m)

	collection := db.Collection("oidcCacheCollection")
	// Insert the item
	_, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}
	return nil
}

// Get gets the items with the given state.
// If none is found, it will return nil.
func Get(m manipulate.Manipulator, state string) (*CacheItem, error) {

	db := manipmongo.GetDatabase(m)

	item := &CacheItem{}
	collection := db.Collection("oidcCacheCollection")
	filter := bson.M{"state": state}
	err := collection.FindOne(context.TODO(), filter).Decode(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Delete deletes the items with the given state.
func Delete(m manipulate.Manipulator, state string) error {

	db := manipmongo.GetDatabase(m)

	collection := db.Collection("oidcCacheCollection")
	filter := bson.M{"state": state}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
