//go:build integration

package testingh

import (
	"context"
	"strings"

	"github.com/google/uuid"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
)

type DBSuite struct {
	ContextSuite

	DBPrefix string
	Store    *store.Client
	Database *store.Database
	cleanUp  func(ctx context.Context)
}

func NewDBSuite(dbPrefix string) DBSuite {
	return DBSuite{DBPrefix: dbPrefix}
}

func (ds *DBSuite) SetupSuite() {
	ds.ContextSuite.SetupSuite()

	db := ds.DBPrefix + strings.ReplaceAll(uuid.New().String(), "-", "")
	ds.T().Logf("database: %s", db)

	ds.Store, ds.cleanUp = PrepareDB(ds.SuiteCtx, ds.T(), db)
	ds.Database = store.NewDatabase(ds.Store)
}

func (ds *DBSuite) TearDownSuite() {
	if f := ds.cleanUp; f != nil {
		f(ds.SuiteCtx)
	}
	ds.ContextSuite.TearDownSuite()
}
