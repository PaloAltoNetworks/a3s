package indexes

import (
	"fmt"
	"strings"

	"github.com/globalsign/mgo"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.uber.org/zap"
)

// Ensure ensures the indexes declared in the specs are aligned.
func Ensure(m manipulate.Manipulator, model elemental.ModelManager, packageName string) (err error) {

	indexes := GetIndexes(packageName, model)

	for ident, mIndexes := range indexes {
		if err = manipmongo.EnsureIndex(m, ident, mIndexes...); err != nil {
			zap.L().Warn("Unable to ensure index", zap.Error(err))
		}
	}

	return nil
}

// GetIndexes returns all the indexes for all the identity in the model
func GetIndexes(packageName string, model elemental.ModelManager) (mIndexes map[elemental.Identity][]mgo.Index) {

	var indexes [][]string

	mIndexes = map[elemental.Identity][]mgo.Index{}

	for _, ident := range model.AllIdentities() {

		if ident.Package != packageName {
			continue
		}

		indexes = model.Indexes(ident)
		if len(indexes) == 0 {
			continue
		}

		iName := "index_" + ident.Name + "_"

		for i := range indexes {

			idx := mgo.Index{}

			piName := iName
			var hashedApplied bool

			for _, name := range indexes[i] {

				if hashedApplied {
					panic("hashed index must not be a compound index")
				}

				switch name {

				case ":shard":
					piName = "shard_" + iName

				case ":unique":
					idx.Unique = true

				default:

					name = strings.ToLower(name)
					if attSpec, ok := model.Identifiable(ident).(elemental.AttributeSpecifiable); ok {
						if bsonName := attSpec.SpecificationForAttribute(name).BSONFieldName; bsonName != "" {
							name = bsonName
						}
					}

					idx.Key = append(idx.Key, name)

					if strings.HasPrefix(name, "$hashed:") {
						hashedApplied = true
					}
				}
			}

			idx.Name = piName + strings.Join(idx.Key, "_")
			idx.Background = true

			mIndexes[ident] = append(mIndexes[ident], idx)
		}
	}

	return mIndexes
}

func coucou() {
	fmt.Println("not tested")
}
