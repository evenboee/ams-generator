package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/evenboee/go-config"

	_ "github.com/lib/pq"
)

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

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	file = flag.String("file", "", "output file")
)

func main() {
	flag.Parse()

	dbConf := config.MustLoadFile[DBConfig](".env")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name)
	conn := require(sql.Open("postgres", dsn))
	guard(conn.Ping())
	defer conn.Close()

	var out *os.File
	if *file != "" {
		out = require(os.OpenFile(*file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644))
	} else {
		out = os.Stdout
	}

	startT := time.Now()

	prev := 0
	for {
		row := conn.QueryRow(`SELECT COUNT(*) from pg_stat_activity;`)
		var count int
		guard(row.Scan(&count))

		if prev != count {
			prev = count
			t := time.Since(startT).Seconds()
			fmt.Fprintf(out, "%f: %d\n", t, count)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
