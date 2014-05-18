package sup

import (
	"database/sql"
	"log"
)

var migrations = []func() error{
	migrateCreateStatus,
}

func migrateDb() {
	// create the table, if needed
	_, err := db.Exec(`create table if not exists migration (
        version int not null
    ) engine=InnoDB`)
	if err != nil {
		panic(err)
	}

	// what version are we at? lock it
	var version int
	err = db.QueryRow("select version from migration").Scan(&version)
	if err == sql.ErrNoRows {
		version = 0
	} else if err != nil {
		panic(err)
	}

	// apply migrations
	for ; version < len(migrations); version++ {
		log.Printf("Migrating DB version %d to %d\n", version, version+1)
		err = migrations[version]()
		if err != nil {
			panic(err)
		}
		if version == 0 {
			_, err = db.Exec("insert into migration(version) values (?)", version+1)
		} else {
			_, err = db.Exec("update migration set version=?", version+1)
		}
	}
}

func migrateCreateStatus() error {
	_, err := db.Exec(`create table if not exists status (
        user varchar(100) not null,
        create_time varchar(100) not null,
        tags varchar(100) not null,
        description varchar(100) not null,
        ip varchar(100) not null,

        primary key (user, create_time)
    ) engine=InnoDB character set=utf8`)
	return err
}
