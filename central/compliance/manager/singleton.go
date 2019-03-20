package manager

import (
	"github.com/stackrox/rox/pkg/sync"

	"github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/compliance/data"
	"github.com/stackrox/rox/central/compliance/standards"
	complianceResultsStore "github.com/stackrox/rox/central/compliance/store"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	nodeStore "github.com/stackrox/rox/central/node/globalstore"
	"github.com/stackrox/rox/central/scrape/factory"
)

var (
	managerInstance     ComplianceManager
	managerInstanceInit sync.Once
)

// Singleton returns the compliance manager singleton instance.
func Singleton() ComplianceManager {
	managerInstanceInit.Do(func() {
		var err error
		managerInstance, err = NewManager(standards.RegistrySingleton(), ScheduleStoreSingleton(), datastore.Singleton(), nodeStore.Singleton(), deploymentDatastore.Singleton(), data.NewDefaultFactory(), factory.Singleton(), complianceResultsStore.Singleton())
		if err != nil {
			log.Panicf("Could not create compliance manager: %v", err)
		}
	})
	return managerInstance
}
