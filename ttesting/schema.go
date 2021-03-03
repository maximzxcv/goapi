package ttesting

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

//TODO: migration scripts to file
var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "Creating table posts",
			Script: `CREATE TABLE public.users
			(
				id uuid NOT NULL,
				name character varying COLLATE pg_catalog."default" NOT NULL,
				password character varying COLLATE pg_catalog."default",
				CONSTRAINT users_pkey PRIMARY KEY (id),
				CONSTRAINT users_name_key UNIQUE (name)
			)`,
		},
		// {
		// 	Version:     2,
		// 	Description: "Adding column body",
		// 	Script:      "ALTER TABLE public.users OWNER to postgres;",
		// },
	}
)

// Migrate creates DB
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
