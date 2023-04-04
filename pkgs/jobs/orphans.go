package jobs

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.uber.org/zap"
)

var (
	orphanCleaningAdjustment = time.Duration(rand.Intn(60)) * time.Second
)

// ScheduleOrphanedObjectsDeleteJob periodically cleans objects that lives in a
// deleted namespace. It takes 2 manipulators. The first one is used to perform
// object cleanup in the database. The other is used to retrieve the list of
// existing namespaces. The job will remove all the identities provided by the
// given model manager that have a package set to the given packageName. If
// package name is set to *, the job will apply to all identities The job will
// run at the defined period + a random duration between 0 and 1 minute.
func ScheduleOrphanedObjectsDeleteJob(
	ctx context.Context,
	nsm manipulate.Manipulator,
	m manipulate.TransactionalManipulator,
	identities []elemental.Identity,
	period time.Duration,
) {

	ticker := time.NewTicker(period + orphanCleaningAdjustment)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			if err := DeleteOrphanedObjects(ctx, nsm, m, identities); err != nil {
				zap.L().Error(
					"Unable to complete job DeleteOrphanedObjects",
					zap.Error(err),
				)
			}

		case <-ctx.Done():
			return
		}
	}
}

// DeleteOrphanedObjects deletes the objects with the given
// identities that are not in any of the given namespaces.
func DeleteOrphanedObjects(
	ctx context.Context,
	apiManipulator manipulate.Manipulator,
	m manipulate.TransactionalManipulator,
	identities []elemental.Identity,
) error {

	os, err := manipulate.Iter(
		ctx,
		apiManipulator,
		manipulate.NewContext(
			ctx,
			manipulate.ContextOptionRecursive(true),
			manipulate.ContextOptionOrder("ID"),
			manipulate.ContextOptionFields([]string{"name"}),
		),
		api.SparseNamespacesList{},
		0,
	)
	if err != nil {
		return fmt.Errorf("unable to retrieve list of namespaces: %w", err)
	}

	namespaces := os.List()
	names := make([]string, len(namespaces)+1)
	names[0] = "/"
	for i, ns := range namespaces {
		names[i+1] = *(ns.(*api.SparseNamespace).Name)
	}

	mgodb, closer, err := manipmongo.GetDatabase(m)
	if err != nil {
		return err
	}
	defer closer()

	var orphans []struct {
		ID interface{} `bson:"_id"`
	}

	for _, i := range identities {

		// Never run against namespaces as it is our source of truth
		if i.Name == api.NamespaceIdentity.Name {
			continue
		}

		if err = mgodb.C(i.Name).Pipe([]bson.M{
			{
				"$match": bson.M{
					"namespace": bson.M{
						"$nin": names,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": "$_id",
				},
			},
		}).AllowDiskUse().All(&orphans); err != nil {
			return fmt.Errorf("unable to retrieve orphans for '%s': %w", i.Category, err)
		}

		if len(orphans) == 0 {
			continue
		}

		zap.L().Debug("Deleting orphans",
			zap.String("identity", i.Name),
			zap.Int("count", len(orphans)),
			zap.Reflect("orphans", orphans),
		)

		ids := make([]any, 0, len(orphans))
		for _, orphan := range orphans {
			ids = append(ids, orphan.ID)
		}

		if _, err = mgodb.C(i.Name).RemoveAll(bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		}); err != nil {
			return fmt.Errorf("unable to remove orphans for '%s': %w", i.Category, err)
		}
	}

	return nil
}
