// +build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	source := pgtest.GetConnectionString(t)
	config, err := pgxpool.ParseConfig(source)
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	Destroy(pool)
	store := New(pool)

	alert := fixtures.GetAlert()
	foundAlert, exists, err := store.Get(alert.GetId())
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Nil(t, foundAlert)

	assert.NoError(t, store.Upsert(alert))
	foundAlert, exists, err = store.Get(alert.GetId())
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, alert, foundAlert)

	alert.State = storage.ViolationState_RESOLVED
	assert.NoError(t, store.Upsert(alert))

	foundAlert, exists, err = store.Get(alert.GetId())
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, alert, foundAlert)

	assert.NoError(t, store.Delete(alert.GetId()))
	foundAlert, exists, err = store.Get(alert.GetId())
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Nil(t, foundAlert)

	alert.State = storage.ViolationState_ACTIVE
	assert.NoError(t, store.Upsert(alert))

	indexer := NewIndexer(pool)

	// Common alert searches
	results, err := indexer.Search(search.NewQueryBuilder().AddExactMatches(search.DeploymentID, alert.GetDeployment().GetId()).ProtoQuery())
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	q := search.NewQueryBuilder().
		AddExactMatches(search.DeploymentID, alert.GetDeployment().GetId()).
		AddExactMatches(search.PolicyID, alert.GetPolicy().GetId()).
		AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String()).
		ProtoQuery()
	results, err = indexer.Search(q)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}
