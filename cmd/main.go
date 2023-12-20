package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

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

var (
	mode       = flag.String("mode", "none", "run async generation")
	configFile = flag.String("config", "config.yaml", "config file")
	envFile    = flag.String("env", ".env", "env file")
	maxConn    = flag.Int("conns", 0, "max connections")
)

func main() {
	flag.Parse()

	if *mode == "none" {
		fmt.Println("specify mode")
		return
	}

	genConf := require(generator.LoadConfig(*configFile))
	dbConf := config.MustLoadFile[DBConfig](*envFile)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name)
	conn := require(sql.Open("postgres", dsn))
	guard(conn.Ping())

	if *maxConn > 0 {
		conn.SetMaxOpenConns(*maxConn)
	}

	store := db.NewStore(conn)

	switch *mode {
	case "async":
		err := generator.GenerateAsync(context.Background(), store, genConf)
		if err != nil {
			panic(err)
		}
	case "sync":
		err := generator.Generate(context.Background(), store, genConf)
		if err != nil {
			panic(err)
		}
	case "3":
		genConf2 := require(generator.LoadConfig2(*configFile))

		errChan := generator.GenerateAsync2(context.Background(), store, genConf2)
		for err := range errChan {
			if err != nil {
				log.Printf("error %T %s\n", err, err.Error())
			}
		}
	default:
		panic("invalid mode")
	}
}
