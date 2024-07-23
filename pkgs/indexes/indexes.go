package indexes

import (
	"strings"

	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func GetIndexes(packageName string, model elemental.ModelManager) (mIndexes map[elemental.Identity][]mongo.IndexModel) {

	var indexes [][]string

	mIndexes = map[elemental.Identity][]mongo.IndexModel{}

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

			piName := iName
			var hashedApplied bool
			var keys []string
			var unique bool

			for _, name := range indexes[i] {

				if hashedApplied {
					panic("hashed index must not be a compound index")
				}

				switch name {

				case ":shard":
					piName = "shard_" + iName

				case ":unique":
					unique = true

				default:

					name = strings.ToLower(name)
					if attSpec, ok := model.Identifiable(ident).(elemental.AttributeSpecifiable); ok {
						if bsonName := attSpec.SpecificationForAttribute(name).BSONFieldName; bsonName != "" {
							name = bsonName
						}
					}

					keys = append(keys, name)

					if strings.HasPrefix(name, "$hashed:") {
						hashedApplied = true
					}
				}
			}

			idxKeys := bson.D{}
			for _, key := range keys {
				idxKeys = append(idxKeys, bson.E{Key: key, Value: 1}) // Use 1 for ascending order
			}

			idx := mongo.IndexModel{
				Keys:    idxKeys,
				Options: options.Index().SetName(piName + strings.Join(keys, "_")).SetUnique(unique),
			}

			mIndexes[ident] = append(mIndexes[ident], idx)
		}
	}

	return mIndexes
}
