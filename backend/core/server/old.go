package server

// // Handles a /search request, queries database for rows
// // containing ALL keys and
// // wrties output to response writer
// func handleSearch(serverParams serverFuncParams, keys []string) {
//
// 	if len(keys) == 0 {
// 		fmt.Fprintf(serverParams.writer, emptyJSON)
// 		return
// 	}
//
// 	// TODO: All this is wrong
//
// 	query := strings.Join(keys, " ")
//
// 	rm := Folder.ReverseMappingLocal()
//
// 	results := search.Search(Config, &Folder, rm, query)
// 	response := utils.SearchResponse{
// 		Query:   query,
// 		Results: results,
// 	}
//
// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		fmt.Fprintf(serverParams.writer, "JSON failed: %s\n", err)
// 		return
// 	}
//
// 	fmt.Fprintf(serverParams.writer, "%s\n", jsonResponse)
//
// }

// // Querys the database for rows containing ALL the given keys.
// // Writes output to writer
// func queryJSONKeysAll(db *sql.DB, writer io.Writer, keys []string) {
// 	query := `SELECT * FROM page WHERE dict ?& $1`
//
// 	if len(keys) == 0 {
// 		panic(`No keys given`)
// 	}
//
// 	fmt.Printf("%s (%s)\n", query, keys)
//
// 	rows, err := db.Query(query, pq.StringArray(keys))
// 	checkSQLError(err)
// 	defer func() {
// 		err = rows.Close()
// 		checkIOError(err)
// 	}()
//
// 	writeRows(writer, rows)
// }
