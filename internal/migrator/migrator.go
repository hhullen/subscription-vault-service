package migrator

import (
	"context"
	"database/sql"
	"log"

	"subscription-vault-service/internal/clients/postgres"
	ds "subscription-vault-service/internal/datastruct"
	secretprovider "subscription-vault-service/internal/secret_provider"
	"subscription-vault-service/internal/supports"

	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations/postgres"
	cmdUp         = "up"
	cmdDown       = "down"
	cmdDownAll    = "down-all"
	cmdStatus     = "status"

	dbDialect = "postgres"

	db_host_secret_path     = "db_host"
	db_port_secret_path     = "db_port"
	db_password_secret_path = "db_migrator_password"
	db_user_secret_path     = "db_migrator_user"
	db_name_secret_path     = "db_name"
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

	secretDir := ds.DefaultSecretsDir
	if supports.IsInContainer() {
		secretDir = ds.DefaultContainerSecretsDir
	}

	sp := secretprovider.NewSecretProvider(secretDir)

	host, err := sp.ReadSecret(db_host_secret_path)
	if err != nil {
		log.Fatalf("readinng db secret: %s", err.Error())
		return
	}
	port, err := sp.ReadSecret(db_port_secret_path)
	if err != nil {
		log.Fatalf("readinng db secret: %s", err.Error())
		return
	}
	user, err := sp.ReadSecret(db_user_secret_path)
	if err != nil {
		log.Fatalf("readinng db secret: %s", err.Error())
		return
	}
	password, err := sp.ReadSecret(db_password_secret_path)
	if err != nil {
		log.Fatalf("readinng db secret: %s", err.Error())
		return
	}
	dbname, err := sp.ReadSecret(db_name_secret_path)
	if err != nil {
		log.Fatalf("readinng db secret: %s", err.Error())
		return
	}

	if !supports.IsInContainer() {
		host = "localhost"
	}
	db, err := postgres.NewSQLConn(ctx, user, password, host, port, dbname)
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
