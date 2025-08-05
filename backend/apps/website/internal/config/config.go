package config

type Config struct {
	App      App
	Database Database
}

type App struct {
	Port         string `doc:"Port to listen on"                                 default:"8000"`
	IsProduction bool   `doc:"Set the logger according to this boolean"          default:"false"`
	Version      string `doc:"Application version. Will be on the documentation" default:"development"`
}

type Database struct {
	DSN string `doc:"DSN of the database" default:"postgres://postgres:@localhost:5432/gradient?sslmode=disable"`
}
