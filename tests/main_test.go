package main

import (
	"testing"

	"gen-test/database/queries"

	"github.com/kr/pretty"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_main(t *testing.T) {
	Convey("Test_relation", t, func() {
		dsn := "host=10.1.1.2 user=postgres password=xixi0202 dbname=test port=5433 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
		So(err, ShouldBeNil)
		queries.SetDefault(db)

		Convey("many 2 many", func() {
			tbl, query := queries.Student.QueryContext(t.Context())
			So(tbl, ShouldNotBeEmpty)

			m, err := query.Preload(tbl.Class).Where(tbl.ID.Eq(1)).First()
			So(err, ShouldBeNil)

			pretty.Print(m)
		})
	})
}
