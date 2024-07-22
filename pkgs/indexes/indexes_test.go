package indexes

import (
	"reflect"
	"testing"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetIndexes(t *testing.T) {
	type args struct {
		packageName string
		identity    elemental.Identity
		model       elemental.ModelManager
	}
	tests := []struct {
		name         string
		args         args
		wantMIndexes map[elemental.Identity][]mongo.IndexModel
	}{
		{
			name: "all indexes",
			args: args{
				packageName: "a3s",
				identity:    api.AuthorizationIdentity,
				model:       api.Manager(),
			},
			wantMIndexes: map[elemental.Identity][]mongo.IndexModel{
				api.AuthorizationIdentity: {
					{
						Keys:    bson.D{{Key: "zone", Value: 1}, {Key: "zhash", Value: 1}},
						Options: options.Index().SetName("shard_index_authorization_zone_zhash").SetBackground(true).SetUnique(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace").SetBackground(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}, {Key: "_id", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace__id").SetBackground(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}, {Key: "flattenedsubject", Value: 1}, {Key: "disabled", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace_flattenedsubject_disabled").SetBackground(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}, {Key: "flattenedsubject", Value: 1}, {Key: "propagate", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace_flattenedsubject_propagate").SetBackground(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}, {Key: "importlabel", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace_importlabel").SetBackground(true),
					},
					{
						Keys:    bson.D{{Key: "namespace", Value: 1}, {Key: "trustedissuers", Value: 1}},
						Options: options.Index().SetName("index_authorization_namespace_trustedissuers").SetBackground(true),
					},
				},
			},
		},
	}

	// compareIndexOptions compares two IndexOptions structs using reflection
	compareIndexOptions := func(expected, actual *options.IndexOptions) bool {
		if expected == nil && actual == nil {
			return true
		}

		if expected == nil || actual == nil {
			return false
		}

		if expected.Name == nil && actual.Name != nil ||
			expected.Name != nil && actual.Name == nil ||
			*expected.Name != *actual.Name {
			return false
		}

		expectedBackground := expected.Background != nil && *expected.Background == true
		actualBackground := actual.Background != nil && *actual.Background == true

		expectedUnique := expected.Unique != nil && *expected.Unique == true
		actualUnique := actual.Unique != nil && *actual.Unique == true

		if expectedUnique != actualUnique || expectedBackground != actualBackground {
			return false
		}
		return true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotMIndexes := GetIndexes(tt.args.packageName, tt.args.model)
			actual := gotMIndexes[tt.args.identity]
			expected := tt.wantMIndexes[tt.args.identity]

			if len(expected) != len(actual) {
				t.Errorf("Number of index models do not match. Expected: %d, got: %d", len(expected), len(actual))
				return
			}

			for i := range expected {
				if !reflect.DeepEqual(expected[i].Keys, actual[i].Keys) {
					t.Errorf("Keys do not match. Expected: %v, got: %v", expected[i].Keys, actual[i].Keys)
				}

				if !compareIndexOptions(expected[i].Options, actual[i].Options) {
					t.Errorf("Options do not match. Expected: %+v, got: %+v", expected[i].Options, actual[i].Options)
				}
			}
		})
	}
}
