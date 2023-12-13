package main

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/evenboee/ams-generator/db/sqlc"
	"github.com/evenboee/ams-generator/generator"
	"github.com/evenboee/go-config"

	_ "github.com/lib/pq"
)

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

type DBConfig struct {
	User string `config:"DB_USER"`
	Pass string `config:"DB_PASS"`
	Name string `config:"DB_NAME"`
	Host string `config:"DB_HOST"`
	Port int    `config:"DB_PORT"`
}

func require[T any](d T, err error) T {
	guard(err)
	return d
}

func main() {
	genConf := require(generator.LoadConfig("config.yaml"))
	dbConf := config.MustLoadFile[DBConfig](".env")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name)
	conn := require(sql.Open("postgres", dsn))
	guard(conn.Ping())
	defer conn.Close()

	store := db.NewStore(conn)

	err := generator.Generate(context.Background(), store, genConf)
	if err != nil {
		panic(err)
	}
}
