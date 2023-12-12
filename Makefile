dbml2sql:
	dbml2sql db.dbml -o db/migration/000001_init.up.sql

# make createMigration NAME={{migration_name}}
createMigration:
	migrate create -ext sql -dir db/migration -seq "$(NAME)"

.PHONY: dbml2sql createMigration
