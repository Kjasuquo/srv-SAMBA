package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	DbMigrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// The SQL code to upgrade to a more recent version of the database schema
		_, err := db.Exec("CREATE TABLE _test(id SERIAL PRIMARY KEY)")
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		// The SQL code to downgrade to the previous version of the database schema
		_, err := db.Exec("DROP TABLE _test")
		return err
	})
}
