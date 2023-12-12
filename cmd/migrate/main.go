package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/evenboee/go-config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	migrationURL = flag.String("migration-url", "file://db/migration", "migration url")
	yes          = flag.Bool("yes", false, "skip confirmation prompt")
)

type Config struct {
	DbUser string `config:"DB_USER" default:"postgres"`
	DbPass string `config:"DB_PASS" default:"postgres"`
	DbHost string `config:"DB_HOST" default:"localhost"`
	DbPort string `config:"DB_PORT" default:"5432"`
	DbName string `config:"DB_NAME" default:"postgres"`
}

func main() {
	flag.Parse()

	config := config.MustLoadFile[Config](".env")

	args := flag.Args()

	if len(args) == 0 || len(args) > 2 {
		log.Fatalf("Wrong number of arguments. Expected 1 or 2 arguments: [up|down|version|clean] [steps] got %d\n", len(args))
	}

	dbSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)

	up := true
	if len(args) > 0 {
		switch args[0] {
		case "up", "down":
			up = args[0] == "up"
		case "version", "v":
			version, dirty, err := getVersion(*migrationURL, dbSource)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("version: %d\ndirty:   %t\n", version, dirty)
			os.Exit(0)
		case "force":
			if len(args) != 2 {
				log.Fatalf("Wrong number of arguments. Expected 2 arguments: force [version] got %d\n", len(args))
			}

			version, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("failed to parse version value %q to integer: %s", args[1], err)
			}

			force(*migrationURL, dbSource, version)
			os.Exit(0)
		case "clean":
			clean(*migrationURL, dbSource)
			os.Exit(0)
		default:
			log.Fatalf("invalid direction: %s", args[0])
		}
	}

	steps := 0
	var err error
	if len(args) > 1 {
		steps, err = strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("failed to parse steps value %q to integer: %s", args[1], err)
		}
	}

	version, _, err := getVersion(*migrationURL, dbSource)
	if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			log.Fatalln(err)
		}
		version = 0
	}

	toVersion := "latest"
	if steps > 0 {
		offset := uint(steps)
		if !up {
			if version < offset {
				log.Fatalf("cannot run down migration %d steps when current version is %d", steps, version)
			}
			offset = -offset
		}

		toVersion = fmt.Sprintf("%d", version+offset)
	} else if steps == 0 {
		if !up {
			toVersion = "0"
		}
	}

	fmt.Printf("Migrating version %d -> %s\n", version, toVersion)

	if !up && !*yes {
		if steps != 0 {
			if !YesNoPrompt(fmt.Sprintf("Are you sure you want to run down migration %d steps?", steps), false) {
				log.Fatalf("aborting")
			}
		} else if !YesNoPrompt("Are you sure you want to run down migration?", false) {
			log.Fatalf("aborting")
		}
	}

	runMigration(*migrationURL, dbSource, steps, up)
}

func clean(migrationURL string, dbSource string) {
	version, dirty, err := getVersion(migrationURL, dbSource)
	if err != nil {
		log.Fatalln(err)
	}

	if dirty {
		migration, err := migrate.New(migrationURL, dbSource)
		if err != nil {
			log.Fatalln("failed to create migration:", err)
		}

		if err := migration.Force(int(version)); err != nil {
			log.Fatalln("failed to force version:", err)
		}
	}
}

func force(migrationURL string, dbSource string, version int) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("failed to create migration:", err)
	}

	if err := migration.Force(version); err != nil {
		log.Fatalln("failed to force version:", err)
	}
}

type Logger struct{}

// Printf is like fmt.Printf
func (l Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Verbose should return true when verbose logging output is wanted
func (l Logger) Verbose() bool {
	return true
}

// https://dev.to/tidalcloud/interactive-cli-prompts-in-go-3bj9
func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func runMigration(migrationURL string, dbSource string, steps int, up bool) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("failed to create migration:", err)
	}

	migration.Log = Logger{}

	if steps != 0 {
		if !up {
			steps = -steps
		}

		err = migration.Steps(steps)
	} else {
		if up {
			err = migration.Up()
		} else {
			err = migration.Down()
		}
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("no migration changes")
		}
		log.Fatalf("failed to migrate: %s (%T)\n", err, err)
	}
}

func getVersion(migrationURL string, dbSource string) (uint, bool, error) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration: %w", err)
	}

	migration.Log = Logger{}

	version, dirty, err := migration.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}
