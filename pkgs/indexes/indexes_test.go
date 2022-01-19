package indexes

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
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
		wantMIndexes map[elemental.Identity][]mgo.Index
	}{
		{
			name: "all indexes",
			args: args{
				packageName: "a3s",
				identity:    api.AuthorizationIdentity,
				model:       api.Manager(),
			},
			wantMIndexes: map[elemental.Identity][]mgo.Index{
				api.AuthorizationIdentity: {
					{
						Name:       "index_authorization_namespace_flattenedsubject_disabled",
						Key:        []string{"namespace", "flattenedsubject", "disabled"},
						Background: true,
					},
					{
						Name:       "index_authorization_namespace_flattenedsubject_propagate",
						Key:        []string{"namespace", "flattenedsubject", "propagate"},
						Background: true,
					},
					{
						Name:       "shard_index_authorization_zone_zhash",
						Key:        []string{"zone", "zhash"},
						Background: true,
						Unique:     true,
					},
					{
						Name:       "index_authorization_namespace",
						Key:        []string{"namespace"},
						Background: true,
					},
					{
						Name:       "index_authorization_namespace__id",
						Key:        []string{"namespace", "_id"},
						Background: true,
					},
					{
						Name:       "index_authorization_namespace_importlabel",
						Key:        []string{"namespace", "importlabel"},
						Background: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotMIndexes := GetIndexes(tt.args.packageName, tt.args.model); !reflect.DeepEqual(gotMIndexes[tt.args.identity], tt.wantMIndexes[tt.args.identity]) {
				t.Errorf("GetIndexes()\n"+
					"EXPECTED:\n"+
					"%+v\n"+
					"ACTUAL:\n"+
					"%+v\n",
					tt.wantMIndexes[tt.args.identity],
					gotMIndexes[tt.args.identity])
			}
		})
	}
}
