package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	sharedDb "github.com/SAMBA-Research/microservice-shared/db"
	//"github.com/SAMBA-Research/microservice-template/cmd/migrations"
	"github.com/kjasuquo/srv-SAMBA/internal/config"
	"github.com/kjasuquo/srv-SAMBA/internal/utils"
)

type DbConnectionManager struct {
	utils.WithRWMutex
	connections map[string]*bun.DB
	cfg         *config.Config
	ctx         context.Context
}

func NewDbConnectionManager(cfg *config.Config, ctx context.Context) (cm *DbConnectionManager) {
	return &DbConnectionManager{
		connections: map[string]*bun.DB{},
		cfg:         cfg,
		ctx:         ctx,
	}
}

// NewDbConnection returns a DB, given a simple tenant name (e.g. "samba" or "shared").
// It does so by caching all the connections and returning an existing connection from
// the cache, if it's there.
// The caller must not close this connection.
func (cm *DbConnectionManager) GetDbConnection(tenant string) (db *bun.DB, err error) {
	cnKey := fmt.Sprintf("%s.%s", cm.cfg.Environment, tenant)
	cm.RWMutex.RLock()
	db, found := cm.connections[cnKey]
	cm.RWMutex.RUnlock()

	if found {
		return db, nil
	}
	cm.RWMutex.Lock()
	// Double check
	db, found = cm.connections[cnKey]
	if !found {
		db, err = sharedDb.NewDbConnection(cnKey)
		if err == nil {
			cm.dbMigrate(db, tenant, cnKey)
			cm.connections[cnKey] = db
		}
	}
	cm.RWMutex.Unlock()
	return db, err
}

func (cm *DbConnectionManager) dbMigrate(db *bun.DB, tenant, cnKey string) {
	migrator := migrate.NewMigrator(db, migrations.DbMigrations)

	if !cm.dbTableExists(db, "bun_migrations", tenant) {
		err := migrator.Init(cm.ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error initing migrations")
		}
	}
	/*
		err := migrator.Lock(cm.ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error locking Bun Migrator")
		}
		defer migrator.Unlock(cm.ctx)
	*/
	group, err := migrator.Migrate(cm.ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error in database migration for tenant %s", cnKey)
	}
	if group.ID != 0 {
		log.Info().Interface("New version", group.ID).Msgf("Migrated tenant %s to a new version", cnKey)
	}
}

func (cm *DbConnectionManager) dbTableExists(db *bun.DB, table, schema string) (exists bool) {
	ctx := context.Background()
	var count int
	err := db.NewRaw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ? AND table_schema = ?", table, schema).Scan(ctx, &count)
	if err != nil || count == 0 {
		return false
	}
	return true
}
