package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/mysql" // MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file"    // File source for migrations
	configuration "github.com/tolubydesign/todo-go/app/config"
)

// Config holds MySQL connection details. This would typically come from environment variables or a config file.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Setup and pass the configuration variable necessary to connect to database
func DatabaseConfig(logging *zap.Logger) (Config, error) {
	// NOTE.Learning that the zap logger is provided through uber-fx. Like magic. And through uber-fx, im able to provide db.Config
	c, err := configuration.GetConfiguration()
	logging.Info("DatabaseConfig called")
	if err != nil {
		logging.Warn("unable to retrieve environment variables", zap.String("error", err.Error()))
		panic(err)
	}

	dbCfg := Config{
		Host:     c.Configuration.Mysql.SqlHost,
		Port:     c.Configuration.Mysql.SqlPort,
		User:     c.Configuration.Mysql.SqlUser,
		Password: c.Configuration.Mysql.SqlPassword,
		Database: c.Configuration.Mysql.SqlDatabase,
	}

	if dbCfg.Host == "" {
		// Provide a fallback or return an error if DSN is mandatory
		return Config{}, errors.New("MYSQL_HOST environment variable not provided")
	}
	// TODO: additional error responses. keep user int the loop

	logging.Info("configuration set")
	return dbCfg, nil
}

func DatabaseSourceName(cfg Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database) // root:yourpassword@tcp(127.0.0.1:3306)/test
	return dsn
}

// MySQLConnection is the Fx Provider function for the database connection.
func MySQLConnection(cfg Config, lc fx.Lifecycle) (*sql.DB, error) {
	log.Println("MySQL Connection")

	dsn := DatabaseSourceName(cfg)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		msg := fmt.Sprintf("failed to open database connection: %s", err.Error())
		return nil, errors.New(msg)
	}

	// Ping the database to ensure connection is live
	if err := db.Ping(); err != nil {
		db.Close() // Close on failed ping
		msg := fmt.Sprintf("failed to open database connection: %s", err.Error())
		return nil, errors.New(msg)
	}

	log.Println("MySQL Connection. database pinged")

	// Fx Lifecycle Hook to gracefully close the connection
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Connecting to MySQL...")
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Closing MySQL connection...")
			return db.Close()
		},
	})
	return db, nil
}

// ToDo is the model used in the service layer for data structuring.
type ToDo struct {
	ID               int
	Task             string
	Task_description string
	Created_at       time.Time
	Due_date         *time.Time
}

// ToDoService defines the business logic methods.
type ToDoService struct {
	db     *sql.DB
	logger *zap.Logger
}

// uber-fx automatically handles resolving sql and zap.
func NewToDoService(db *sql.DB, logger *zap.Logger) *ToDoService {
	return &ToDoService{
		db:     db,
		logger: logger.Named("ToDoService"),
	}
}

// CreateToDo adds new todo with the provided details.
func (s *ToDoService) GetToDo(ctx context.Context, limit string, page string) ([]ToDo, error) {
	s.logger.Info("db get todo")
	const query = "SELECT * FROM todo ORDER BY id ASC LIMIT ? OFFSET ?"
	var todos []ToDo

	// Execute the query
	rows, err := s.db.QueryContext(ctx, query, limit, page)
	if err != nil {
		s.logger.Error("Failed to get ToDos", zap.Error(err),
			zap.String("limit", limit), zap.String("page", page))
		return nil, fmt.Errorf("database execution error: %w", err)
	}
	s.logger.Info("db get todo - close()")
	defer rows.Close()

	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Task_description, &todo.Created_at, &todo.Due_date); err != nil {
			s.logger.Error("Failed to Scan Todo", zap.Error(err), zap.Int("id", todo.ID))
			return nil, fmt.Errorf("todo %v: %v", todo.ID, err)
		}
		todos = append(todos, todo)
	}

	s.logger.Info("ToDo got successfully via SQL", zap.String("limit", limit), zap.String("page", page))
	return todos, nil
}

// CreateToDo adds new todo with the provided details.
func (s *ToDoService) CreateToDo(ctx context.Context, task string, description string, due_date *time.Time) (*ToDo, error) {
	// const insertSQL = "INSERT INTO todo (task, task_description, due_date) VALUES (?, ?, ?)"
	params := []string{}
	args := []any{}
	marks := []string{}

	// add task param
	params = append(params, "task")
	args = append(args, task)

	// add description param
	params = append(params, "task_description")
	args = append(args, description)

	// if due_date != nil && *due_date != "" {
	if due_date != nil {
		params = append(params, "due_date")
		args = append(args, *due_date)
	}

	for i := 0; i < len(args); i++ {
		marks = append(marks, "?")
	}

	query := fmt.Sprintf("INSERT INTO todo (%s) VALUES (%s)", strings.Join(params, ", "), strings.Join(marks, ", "))

	s.logger.Info("request structure", zap.String("query", query))
	var date_time time.Time
	if due_date != nil {
		date_time = *due_date
	}

	// Execute the query
	// result, err := s.db.ExecContext(ctx, insertSQL, task, description)
	// if err != nil {
	// 	s.logger.Error("Failed to create ToDo item using SQL Exec", zap.Error(err),
	// 		zap.String("task", task), zap.String("description", description))
	// 	return nil, fmt.Errorf("database execution error: %w", err)
	// }

	// Get the ID of the newly inserted record
	// lastID, err := result.LastInsertId()
	// if err != nil {
	// 	s.logger.Warn("Could not retrieve LastInsertId", zap.Error(err))
	// }

	// Construct the ToDo object with the assigned ID
	t := &ToDo{
		ID: 0,
		// ID:               int(lastID),
		Task:             task,
		Task_description: description,
		Due_date:         &date_time,
		// TODO: return created at date
		// Created_at: ,
	}

	// s.logger.Info("ToDo created successfully via SQL", zap.String("task", task), zap.Int64("id", lastID))
	s.logger.Info("ToDo created successfully via SQL", zap.String("task", task), zap.Int64("id", 0))
	return t, nil
}

// Update existing todo in the database, with the provided todo details.
func (s *ToDoService) UpdateToDo(ctx context.Context, id *int, task *string, description *string, due_date *string) error {
	var err error
	// End result should be "UPDATE todo SET task = ?, task_description = ? WHERE id = ?" if all function parameters are provided
	params := []string{}
	args := []any{}

	if task != nil && *task != "" {
		params = append(params, "task = ?")
		args = append(args, *task)
	}

	if description != nil && *description != "" {
		params = append(params, "task_description = ?")
		args = append(args, *description)
	}

	if due_date != nil && *due_date != "" {
		params = append(params, "due_date = ?")
		args = append(args, *due_date)
	}

	if len(params) == 0 {
		s.logger.Warn("insufficient parameters")
		return fmt.Errorf("insufficient parameters provided")
	}

	updateQuery := fmt.Sprintf("UPDATE todo SET %s WHERE id = ?", strings.Join(params, ", "))
	args = append(args, id)
	s.logger.Info("UPDATE SERVICE:", zap.String("query", updateQuery))
	s.logger.Info("UPDATE SERVICE: arguments", zap.Any("args", args))

	// Execute the query
	_, err = s.db.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		s.logger.Warn("Failed to update ToDo item using SQL Exec", zap.Error(err),
			zap.String("task", *task), zap.String("description", *description), zap.Int("id", *id))
		return fmt.Errorf("database execution error: %w", err)
	}

	return nil
}

// RemoveToDo removes task that match provided id and task from SQL database.
func (s *ToDoService) RemoveToDo(ctx context.Context, id int) error {
	const removeSQL = "DELETE FROM todo WHERE id = ?"

	// Execute the query
	result, err := s.db.ExecContext(ctx, removeSQL, id)
	if err != nil {
		s.logger.Error("Failed to create ToDo item using SQL Exec", zap.Error(err),
			zap.Int("id", id))
		return fmt.Errorf("database execution error: %w", err)
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected > 0 {
		s.logger.Info("SQL DELETE: item delete", zap.Int("id", id))
	}
	return nil
}

// NewDB creates and returns the *sql.DB instance.
// Using fx.Lifecycle to manage things .
func NewDB(lc fx.Lifecycle, cfg Config, logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Attempting to connect to MySQL database using database/sql...")
	dsn := DatabaseSourceName(cfg)

	// Open the standard database connection.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database handle: %w", err)
	}

	// Set connection pool settings (best practice)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Ping the database to verify the connection is active
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to MySQL database.")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// NOTE. uncomment if you struggle to populate data through migration up
			// // A simple CREATE TABLE statement to ensure the necessary table exists.
			// const createTableSQL = `
			// 	CREATE TABLE IF NOT EXISTS todo (
			// 		id INT AUTO_INCREMENT PRIMARY KEY,
			// 		task VARCHAR(255) NOT NULL,
			// 		task_description TEXT,
			// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			// 		due_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			// 	);`

			// logger.Info("Attempting to create table todo")
			// _, err := db.ExecContext(ctx, createTableSQL)
			// if err != nil {
			// 	return fmt.Errorf("failed to create products table: %w", err)
			// }
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection.")
			return db.Close()
		},
	})

	return db, nil
}

// Export the NewToDoService constructor for FX.
var ServiceModule = fx.Options(
	fx.Provide(NewToDoService),
)

// DatabaseModule exports the NewDB constructor for FX.
var DatabaseModule = fx.Options(
	fx.Provide(NewDB),
)
