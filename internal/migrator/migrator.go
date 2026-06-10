package migrator

import (
	"context"
	"database/sql"
	"log"

	// "subscription-vault-service/internal/clients/mysql"
	"subscription-vault-service/internal/clients/postgres"
	secretprovider "subscription-vault-service/internal/secret_provider"
	"subscription-vault-service/internal/supports"

	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations/mysql"
	cmdUp         = "up"
	cmdDown       = "down"
	cmdDownAll    = "down-all"
	cmdStatus     = "status"

	dbDialect = "mysql"

	defaultSecretsDir          = "./secrets/"
	defaultContainerSecretsDir = "/run/secrets/"
)

var executors = map[string]func(db *sql.DB, dir string, opts ...goose.OptionsFunc) error{
	cmdUp:      goose.Up,
	cmdDown:    goose.Down,
	cmdStatus:  goose.Status,
	cmdDownAll: downAll,
}

func Migrate(command string) {
	ctx := context.Background()
	defer ctx.Done()

	secretDir := defaultSecretsDir
	if supports.IsInContainer() {
		secretDir = defaultContainerSecretsDir
	}

	sp := secretprovider.NewSecretProvider(secretDir)

	db, err := postgres.NewSQLConn(ctx, sp)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	ExecMigration(db, command, migrationsDir)

}

func ExecMigration(db *sql.DB, command, migrationsDir string) {
	executor, exists := executors[command]
	if !exists {
		log.Fatalf("Wrong comand send: %s. Required: %s/%s/%s/%s", command, cmdUp, cmdDown, cmdDownAll, cmdStatus)
	}

	goose.SetBaseFS(nil)
	if err := goose.SetDialect(dbDialect); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := executor(db, migrationsDir); err != nil {
		log.Fatalf("failed applying migrations: %v", err)
	}

	log.Printf("%s successfully migrated\n", migrationsDir)
}

func downAll(db *sql.DB, dir string, opts ...goose.OptionsFunc) error {
	return goose.DownTo(db, dir, 0, opts...)
}
