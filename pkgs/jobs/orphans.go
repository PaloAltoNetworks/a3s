package jobs

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
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
			manipulate.ContextOptionFields([]string{"namespace", "deleteTime"}),
		),
		api.SparseNamespaceDeletionRecordsList{},
		0,
	)
	if err != nil {
		return fmt.Errorf("unable to retrieve list of namespaces: %w", err)
	}

	deletionRecords := os.List()
	if len(deletionRecords) == 0 {
		return nil
	}

	for _, deletionRecord := range deletionRecords {

		namespace := *(deletionRecord.(*api.SparseNamespaceDeletionRecord).Namespace)
		deletionDate := *(deletionRecord.(*api.SparseNamespaceDeletionRecord).DeleteTime)

		mctx := manipulate.NewContext(
			ctx,
			manipulate.ContextOptionFilter(
				elemental.NewFilterComposer().And(
					manipulate.NewNamespaceFilter(namespace, true),
					elemental.NewFilterComposer().Or(
						elemental.NewFilterComposer().WithKey("createTime").NotExists().Done(),
						elemental.NewFilterComposer().WithKey("createTime").LesserThan(deletionDate).Done(),
					).Done(),
				).Done(),
			),
		)

		for _, i := range identities {

			if i.IsEqual(api.NamespaceDeletionRecordIdentity) {
				continue
			}

			if err := m.DeleteMany(mctx.Derive(), i); err != nil {
				return fmt.Errorf("unable to deletemany '%s' in namespace '%s': %w", i.Category, namespace, err)
			}
		}
	}

	return nil
}
