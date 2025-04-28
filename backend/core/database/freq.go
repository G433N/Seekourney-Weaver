package database

import (
	"database/sql"
	"log"
	"seekourney/utils"

	"github.com/lib/pq"
)

type sqlResult struct {
	path  utils.Path
	score utils.Frequency
}

func (sqlPath sqlResult) SQLScan(rows *sql.Rows) (sqlResult, error) {
	var path utils.Path
	var score utils.Frequency
	err := rows.Scan(&path, &score)
	if err != nil {
		return sqlResult{}, err
	}
	return sqlResult{
		path:  path,
		score: score,
	}, nil
}

func (sqlPath sqlResult) IntoKey() string {
	return string(sqlPath.path)
}

func (sqlPath sqlResult) IntoValue() utils.Frequency {
	return sqlPath.score
}

func FreqMap(db *sql.DB, word utils.Word) (utils.WordFrequencyMap, error) {

	wordStr := string(word)

	json := JsonValue("words", wordStr, "score")
	q := Select().Queries("path", json).From("document").Where("words ?& $1")

	w := []string{wordStr}
	rows, err := db.Query(string(q), pq.StringArray(w))

	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
			panic(err)
		}
	}()

	rMap, err := ScanRowsIntoMapRaw[sqlResult](rows, func(k string) utils.Path {
		return utils.Path(k)
	}, func(v utils.Frequency) utils.Frequency {
		return v
	})

	if err != nil {
		return nil, err
	}

	return rMap, nil
}
