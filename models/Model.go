package models

import (
	"os"
	"songLibrary/customDb"
	"songLibrary/customLog"
	"songLibrary/utils"
	"strings"
)

type Model struct {
	table  string
	fields map[string]string
}

func (model *Model) Table() string {
	return model.table
}

func (model *Model) SetTable(tableTitle string) {
	model.table = tableTitle
}

func (model *Model) CheckModelTable() bool {
	var resp bool
	if model.table != "" {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)

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

func (model *Model) RunTableMigration() bool {
	var resp bool
	if !model.CheckModelTable() {
		query := model.loadSQLFile(utils.ConcatSlice([]string{model.table, "_up.sql"}))
		if query != "" {
			db := customDb.GetConnect()
			defer customDb.CloseConnect(db)
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
