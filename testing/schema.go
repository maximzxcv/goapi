package testing

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
			Description: "Creating table users",
			Script: `CREATE TABLE public.users
			(
				id uuid NOT NULL,
				name character varying COLLATE pg_catalog."default" NOT NULL,
				password character varying COLLATE pg_catalog."default",
				CONSTRAINT users_pkey PRIMARY KEY (id),
				CONSTRAINT users_name_key UNIQUE (name)
			);
			CREATE TABLE public.calls
			(
				id character varying COLLATE pg_catalog."default" NOT NULL,
				method character varying COLLATE pg_catalog."default" NOT NULL,
				path character varying COLLATE pg_catalog."default" NOT NULL,
				"userId" uuid,
				CONSTRAINT calls_pkey PRIMARY KEY (id),
				CONSTRAINT "callToUser" FOREIGN KEY ("userId")
					REFERENCES public.users (id) MATCH SIMPLE
					ON UPDATE NO ACTION
					ON DELETE CASCADE
					NOT VALID
			);`,
		},
		// {
		// 	Version:     1,
		// 	Description: "Creating table calls",
		// 	Script: `CREATE TABLE public.calls
		// 	(
		// 		id character varying COLLATE pg_catalog."default" NOT NULL,
		// 		method character varying COLLATE pg_catalog."default" NOT NULL,
		// 		path character varying COLLATE pg_catalog."default" NOT NULL,
		// 		"userId" uuid,
		// 		CONSTRAINT calls_pkey PRIMARY KEY (id),
		// 		CONSTRAINT "callToUser" FOREIGN KEY ("userId")
		// 			REFERENCES public.users (id) MATCH SIMPLE
		// 			ON UPDATE NO ACTION
		// 			ON DELETE CASCADE
		// 			NOT VALID
		// 	);`,
		// },
	}
)

// Migrate creates DB
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
