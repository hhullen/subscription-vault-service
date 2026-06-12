package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"subscription-vault-service/internal/api/v1"
	"subscription-vault-service/internal/clients/postgres"
	ds "subscription-vault-service/internal/datastruct"
	gracefulterminator "subscription-vault-service/internal/graceful_terminator"
	"subscription-vault-service/internal/logger"
	secretprovider "subscription-vault-service/internal/secret_provider"
	"subscription-vault-service/internal/service"
	"subscription-vault-service/internal/supports"
)

const (
	address = ":8080"

	db_host_secret_path     = "db_host"
	db_port_secret_path     = "db_port"
	db_password_secret_path = "db_app_password"
	db_user_secret_path     = "db_app_user"
	db_name_secret_path     = "db_name"
)

// @title           Subscription vault service
// @version         1.0
// @description     Service for managing subscriptions
// @termsOfService  http://swagger.io/terms/

// @contact.name   Maksim
// @contact.url    https://github.com/hhullen
// @contact.email  hhullen@gmail.com

// @license.name  Creative Commons Attribution-NonCommercial 4.0 International Public License
// @license.url   https://creativecommons.org/licenses/by-nc/4.0/deed.en

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	apiLog := logger.NewLogger(os.Stdout, "API")
	serviceLog := logger.NewLogger(os.Stdout, "SERVICE")
	dbLog := logger.NewLogger(os.Stdout, "DB")
	gracefulterminator.Add(func() {
		apiLog.Stop()
		serviceLog.Stop()
		dbLog.Stop()
	})

	secretDir := ds.DefaultSecretsDir
	if supports.IsInContainer() {
		secretDir = ds.DefaultContainerSecretsDir
	}

	sp := secretprovider.NewSecretProvider(secretDir)

	host, err := sp.ReadSecret(db_host_secret_path)
	if err != nil {
		dbLog.FatalKV("readinng db secret", "error", err.Error())
		return
	}
	port, err := sp.ReadSecret(db_port_secret_path)
	if err != nil {
		dbLog.FatalKV("readinng db secret", "error", err.Error())
		return
	}
	user, err := sp.ReadSecret(db_user_secret_path)
	if err != nil {
		dbLog.FatalKV("readinng db secret", "error", err.Error())
		return
	}
	password, err := sp.ReadSecret(db_password_secret_path)
	if err != nil {
		dbLog.FatalKV("readinng db secret", "error", err.Error())
		return
	}
	dbname, err := sp.ReadSecret(db_name_secret_path)
	if err != nil {
		dbLog.FatalKV("readinng db secret", "error", err.Error())
		return
	}

	if !supports.IsInContainer() {
		host = "localhost"
	}

	dbConn, err := postgres.NewSQLConn(ctx, user, password, host, port, dbname)
	if err != nil {
		dbLog.FatalKV("connecting db", "error", err.Error())
		return
	}

	gracefulterminator.Add(func() {
		if err := dbConn.Close(); err != nil {
			dbLog.ErrorKV("closing db", "error", err.Error())
		}
	})

	db := postgres.NewClient(ctx, dbConn)

	subsService := service.NewService(ctx, db, serviceLog)

	apiService, err := api.NewAPI(ctx, address, subsService, apiLog)
	if err != nil {
		apiLog.FatalKV("creating api", "error", err.Error())
		return
	}

	gracefulterminator.Add(func() {
		if err := apiService.Stop(); err != nil {
			serviceLog.ErrorKV("stopping API", "error", err.Error())
		}
	})

	go func() {
		if err := apiService.StartListening(); err != nil {
			apiLog.FatalKV("rinning api", "error", err.Error())
			return
		}
	}()

	<-ctx.Done()
	gracefulterminator.Stop()
}
