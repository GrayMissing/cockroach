package sql_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/testutils/serverutils"
	"github.com/cockroachdb/cockroach/util/leaktest"
)

func TestGeoNoIndex(t *testing.T) {
	defer leaktest.AfterTest(t)()

	params, _ := createTestServerParams()
	s, sqlDB, _ := serverutils.StartServer(t, params)
	defer s.Stopper().Stop()

	if _, err := sqlDB.Exec(`
		CREATE DATABASE d;
		SET DATABASE = d;
		CREATE TABLE t (g GEOGRAPHY, INDEX (g));
	`); err != nil {
		t.Fatal(err)
	}

	rnd := rand.New(rand.NewSource(0))

	for j := 0; j < 1; j++ {
		buf := bytes.NewBufferString("INSERT INTO t VALUES\n")
		for i := 0; i < 1; i++ {
			lat := rnd.Float64()*180 - 90
			lng := rnd.Float64()*180 - 90
			if i > 0 {
				buf.WriteString(",\n")
			}
			fmt.Fprintf(buf, "\t"+`('{"type":"Point","coordinates":[%v, %v]}')`, lat, lng)
		}
		if _, err := sqlDB.Exec(buf.String()); err != nil {
			t.Fatal(err)
		}
	}

	var count int
	/*
		if err := sqlDB.QueryRow("SELECT COUNT(*) FROM t").Scan(&count); err != nil {
			t.Fatal(err)
		}
		fmt.Println("COUNT", count)
		//*/

	lat := rnd.Float64()*180 - 90
	lng := rnd.Float64()*180 - 90
	pt := fmt.Sprintf(`'{"type":"Point","coordinates":[%v, %v]}'`, lat, lng)
	q := fmt.Sprintf(`SELECT COUNT(*) FROM t WHERE ST_DISTANCE(g, %s) < 4`, pt)
	/*
		var e1, e2, e3 string
		if err := sqlDB.QueryRow("explain "+q).Scan(&e1, &e2, &e3); err != nil {
			t.Fatal(err)
		}
		fmt.Println("EXPLAIN", e1, e2, e3)
	*/
	fmt.Println(q)
	start := time.Now()
	fmt.Println("BEFORE")
	if err := sqlDB.QueryRow(q).Scan(&count); err != nil {
		t.Fatal(err)
	}
	fmt.Println("AFTER")
	end := time.Since(start)
	fmt.Println("DUR", end)
	fmt.Println("RESULT COUNT", count)
	if count == 0 {
		t.Fatalf("count == 0")
	}
	/*
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		count = 0
		for rows.Next() {
			count++
			var g interface{}
			if err := rows.Scan(&g); err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%T: %v\n", g, g)
		}
		fmt.Println("RESULT COUNT", count)
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	//*/
}
