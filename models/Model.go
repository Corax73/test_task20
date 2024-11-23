package models

import (
	"songLibrary/customDb"
	"songLibrary/customLog"
	"songLibrary/utils"
)

type Model struct {
	table string
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
