package migration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func MigrationsUp(db *sqlx.DB, dsn string) error {
	const op = "internal.storage.postgresql.migration.MigrationsUp()"

	log.Info().Msg("Launching migrations up")
	if db == nil {
		return fmt.Errorf("func:%s  error:%s", op, "database isn`t established")
	}

	migrationDB, err := connectionForMigration(dsn)
	if err != nil {
		return err
	}

	sqlDatabase := migrationDB.DB
	driver, err := newMigrationDriver(sqlDatabase)
	if err != nil {
		return err
	}

	defer closeMigration(driver, migrationDB, op)

	migration, err := newMigrationInstance(driver)
	if err != nil {
		return err
	}

	err = startMigrationUp(migration)
	if err != nil {
		return err
	}

	log.Info().Msg("migrations up successfully")
	return nil
}

func MigrationsDown(db *sqlx.DB, dsn string) error {
	const op = "internal.storage.postgresql.migration.MigrationsDown()"

	log.Info().Msg("Launching migrations down")

	if db == nil {
		return fmt.Errorf("func:%s  error:%s", op, "database isn`t established")
	}

	migrationDB, err := connectionForMigration(dsn)
	if err != nil {
		return err
	}

	sqlDatabase := migrationDB.DB
	driver, err := newMigrationDriver(sqlDatabase)
	if err != nil {
		return err
	}

	defer closeMigration(driver, migrationDB, op)

	migration, err := newMigrationInstance(driver)
	if err != nil {
		return err
	}

	err = startMigrationDown(migration)
	if err != nil {
		return err
	}

	log.Info().Msg("migrations down successfully")
	return nil
}

func connectionForMigration(dsn string) (*sqlx.DB, error) {
	const op = "internal.storage.postgresql.migration.connectionForMigration()"

	log.Debug().Msg("Creating a migration connection")
	migration, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("func:%s  error:%s(%s)", op, err, "failed to create connection for migrations")
	}

	log.Debug().Msg("Migration connection is successful")
	return migration, nil
}

func newMigrationDriver(db *sql.DB) (database.Driver, error) {
	const op = "internal.storage.postgresql.migration.newMigrationDriver()"

	log.Debug().Msg("Creating driver for migrations")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("func:%s  error:%s(%s)", op, err, "couldn't create driver")
	}

	log.Debug().Msg("Migration driver creation is successful")
	return driver, nil
}

func closeMigration(driver database.Driver, migrationDB *sqlx.DB, op string) {
	op += "closeMigration()"
	if err := driver.Close(); err != nil {
		log.Warn().Msgf("func:%s  error:%s", op, "migration's driver couldn't close")
	}

	if err := migrationDB.Close(); err != nil {
		log.Warn().Msgf("func:%s  error:%s", op, "migration's driver couldn't close")
	}
}

func newMigrationInstance(driver database.Driver) (*migrate.Migrate, error) {
	const op = "internal.storage.postgresql.migration.newMigrationInstance()"

	log.Debug().Msg("Creating a migration instance")

	migrationExmpl, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("func:%s  error:%s(%s)", op, err, "coudn't create migrate instance")
	}

	log.Debug().Msg("Migration instance creation is successful")

	return migrationExmpl, nil
}

func startMigrationUp(migration *migrate.Migrate) error {
	const op = "internal.storage.postgresql.migration.startMigrationUp()"

	log.Debug().Msg("Attempting to migration up")

	err := migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("func:%s  error:%s(%s)", op, err, "failed to apply migrations")
	}

	if err == migrate.ErrNoChange {
		log.Debug().Msgf("func:%s  error:%s", op, "no migrations to apply")
	} else {
		log.Debug().Msgf("func:%s  error:%s", op, "migrations applied successfully")
	}
	return nil
}

func startMigrationDown(migration *migrate.Migrate) error {
	const op = "internal.storage.postgresql.migration.startMigrationDown()"
	err := migration.Down()

	log.Debug().Msg("Attempting to migration down")

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("func:%s  error:%s(%s)", op, err, "failed to apply migrations")
	}

	if err == migrate.ErrNoChange {
		log.Debug().Msgf("func:%s  error:%s", op, "no migrations to apply")
	} else {
		log.Debug().Msgf("func:%s  error:%s", op, "migrations applied successfully")
	}
	return nil
}
