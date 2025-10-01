/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/tolubydesign/todo-go/app/db"
	"github.com/tolubydesign/todo-go/app/handler"
	"github.com/tolubydesign/todo-go/app/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle // The fx lifecycle hook
	Mux       *http.ServeMux
}

// NOTE. pressed for time. Will duplicate. Would prefer a more elegant solution.

func MigrateDatabaseUp(p Params, mux *http.ServeMux, logging *zap.Logger, service *db.ToDoService, shutdowner fx.Shutdowner, cfg db.Config, database *sql.DB) {
	logging.Info("Doing an up migration")

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			dsn := db.DatabaseSourceName(cfg)
			driver, _ := mysql.WithInstance(database, &mysql.Config{})

			m, err := migrate.NewWithDatabaseInstance(
				"file://migrate",
				dsn, // "mysql://user:password@tcp(127.0.0.1:3306)/database_name?multiStatements=true",
				driver,
			)

			if err != nil {
				return fmt.Errorf("failed to create UP migrate instance: %w", err)
			}

			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to run UP migrations: %w", err)
			}

			logging.Info("Database UP migrations applied successfully.")
			logging.Info("Lifecycle UP shutdown")
			shutdowner.Shutdown()
			return nil
		},
	})
}

func MigrateDatabaseDown(p Params, mux *http.ServeMux, logging *zap.Logger, service *db.ToDoService, shutdowner fx.Shutdowner, cfg db.Config, database *sql.DB) {
	logging.Info("Doing an down migration")

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			dsn := db.DatabaseSourceName(cfg)
			driver, _ := mysql.WithInstance(database, &mysql.Config{})

			m, err := migrate.NewWithDatabaseInstance(
				"file://migrate",
				dsn,
				driver,
			)

			if err != nil {
				return fmt.Errorf("failed to create DOWN migrate instance: %w", err)
			}

			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to run DOWN migrations: %w", err)
			}

			logging.Info("Database DOWN migrations applied successfully.")
			logging.Info("Lifecycle DOWN shutdown")
			shutdowner.Shutdown()
			return nil
		},
	})
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration command to up-migrate and down-migrate",
	Long:  "Migrate the mysql database up (by adding information) or down (deleting table). The command only accepts the `up` or `down` args",
	Run: func(cmd *cobra.Command, args []string) {
		var app *fx.App

		fmt.Println("migrate called")
		if len(args) == 0 {
			fmt.Println("not enough")
			return
		}

		if len(args) >= 1 {
			fmt.Println("Argument of", args[0], "provided.")

			switch args[0] {
			case "up":

				app = fx.New(
					fx.Provide(
						logging.ZapLogger,
					),
					fx.Provide(
						db.DatabaseConfig,
					),
					db.DatabaseModule,
					db.ServiceModule,
					fx.Provide(
						handler.NewHandler,
						handler.ProvideMux,
					),
					fx.Invoke(MigrateDatabaseUp),
				)

			case "down":
				app = fx.New(
					fx.Provide(
						logging.ZapLogger,
					),
					fx.Provide(
						db.DatabaseConfig,
					),
					db.DatabaseModule,
					db.ServiceModule,
					fx.Provide(
						handler.NewHandler,
						handler.ProvideMux,
					),
					fx.Invoke(MigrateDatabaseDown),
				)

			default:
				fmt.Println("Could not handle request of:", args[0])
				fmt.Println("Only accepting argument of `up` or `down`.")
			}
		}

		if app != nil {
			// starting the application
			app.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
