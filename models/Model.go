package models

import (
	"database/sql"
	"os"
	"slices"
	"songLibrary/customDb"
	"songLibrary/customLog"
	"songLibrary/utils"
	"strings"
)

type Model struct {
	table  string
	Fields map[string]string
}

func (model *Model) Table() string {
	return model.table
}

func (model *Model) SetTable(tableTitle string) {
	model.table = tableTitle
}

func (model *Model) CheckModelTable(db *sql.DB) bool {
	var resp bool
	if model.Table() != "" {
		queryStr := utils.ConcatSlice([]string{
			"SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = '",
			model.Table(),
			"');"})
		rows, err := db.Query(queryStr)
		if err != nil {
			customLog.Logging(err)
		} else {
			for rows.Next() {
				err := rows.Scan(&resp)
				if err != nil {
					customLog.Logging(err)
				}
			}
		}
	}
	return resp
}

func (model *Model) RunTableMigration(db *sql.DB) bool {
	var resp bool
	if !model.CheckModelTable(db) {
		query := model.loadSQLFile(utils.ConcatSlice([]string{model.table, "_up.sql"}))
		if query != "" {
			tx, err := db.Begin()
			if err != nil {
				customLog.Logging(err)
			} else {
				check := true
				for _, q := range strings.Split(string(query), ";") {
					q := strings.TrimSpace(q)
					if q == "" {
						continue
					}
					if _, err := tx.Exec(q); err != nil {
						customLog.Logging(err)
						tx.Rollback()
						check = false
						break
					}
				}
				if check {
					resp = true
					tx.Commit()
				}
			}
		}
	}
	return resp
}

func (model *Model) loadSQLFile(fileName string) string {
	var resp string
	if fileName != "" {
		fileName = utils.ConcatSlice([]string{"./migrations/", fileName})
		_, err := os.Stat(fileName)
		if !os.IsNotExist(err) {
			file, err := (os.ReadFile(fileName))
			if err != nil {
				customLog.Logging(err)
			} else {
				resp = string(file)
			}
		} else {
			customLog.Logging(err)
		}
	}
	return resp
}

func (model *Model) Save() bool {
	var resp bool
	if len(model.Fields) > 0 {
		strSlice := make([]string, 4+((len(model.Fields)-1)*2))
		strSlice = append(strSlice, "INSERT INTO ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " (")
		fields := utils.GetMapKeys(model.Fields)
		index := utils.GetIndexByStrValue(fields, "id")
		fields = slices.Delete(fields, index, index+1)
		strSlice = append(strSlice, strings.Trim(strings.Join(fields, ","), ","))
		strSlice = append(strSlice, ") VALUES (")
		values := utils.GetMapValues(model.Fields)
		valuesToDb := make([]string, len(values))
		for _, val := range values {
			valuesToDb = append(valuesToDb, utils.ConcatSlice([]string{"'", val, "'"}))
		}
		strSlice = append(strSlice, strings.Trim(strings.Join(valuesToDb, ","), ","))
		strSlice = append(strSlice, ") RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		var id int
		err := db.QueryRow(queryStr).Scan(&id)
		if err != nil {
			customLog.Logging(err)
		} else {
			resp = true
		}
	}
	return resp
}
