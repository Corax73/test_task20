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
	table        string
	Fields       map[string]string
	FieldsStruct any
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

func (model *Model) Create(fields map[string]string) map[string]string {
	response := map[string]string{}
	if utils.CompareMapsByStringKeys(model.Fields, fields) {
		model.Fields = fields
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		response = model.Save()
	}
	return response
}

func (model *Model) Save() map[string]string {
	response := map[string]string{}
	if len(model.Fields) > 0 {
		strSlice := make([]string, 5+((len(model.Fields)-1)*2))
		strSlice = append(strSlice, "INSERT INTO ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " (")
		fields := utils.GetMapKeysWithValue(model.Fields)
		index := utils.GetIndexByStrValue(fields, "id")
		if index != -1 {
			fields = slices.Delete(fields, index, index+1)
		}
		strSlice = append(strSlice, strings.Trim(strings.Join(fields, ","), ","))
		strSlice = append(strSlice, ") VALUES (")
		values := utils.GetMapValues(model.Fields)
		valuesToDb := make([]string, len(values))
		for _, val := range fields {
			if _, ok := model.Fields[val]; ok {
				value := model.Fields[val]
				if strings.Contains(value, "'") {
					value = strings.Replace(value, "'", "''", -1)
				}
				valuesToDb = append(valuesToDb, utils.ConcatSlice([]string{"'", value, "'"}))
			}
		}
		strSlice = append(strSlice, strings.Trim(strings.Join(valuesToDb, ","), ","))
		strSlice = append(strSlice, ") RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		var id string
		err := db.QueryRow(queryStr).Scan(&id)
		if err != nil {
			customLog.Logging(err)
		} else {
			response = map[string]string{"id": id}
		}
	}
	return response
}

func (model *Model) Update(fields map[string]string, id string) map[string]string {
	response := map[string]string{}
	fields = utils.GetMapWithoutKeys(fields, []string{"id"})
	if utils.PresenceMapKeysInOtherMap(fields, model.Fields) {
		strSlice := make([]string, 7+((len(fields)-1)*2))
		strSlice = append(strSlice, "UPDATE ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " SET ")
		if len(fields) > 1 {
			strSlice = append(strSlice, "(")
		}
		columns := utils.GetMapKeysWithValue(fields)
		index := utils.GetIndexByStrValue(columns, "id")
		if index != -1 {
			columns = slices.Delete(columns, index, index+1)
		}
		index = utils.GetIndexByStrValue(columns, "title")
		if index != -1 {
			columns = slices.Delete(columns, index, index+1)
		}
		index = utils.GetIndexByStrValue(columns, "group_id")
		if index != -1 {
			columns = slices.Delete(columns, index, index+1)
		}
		if len(fields) > 1 {
			strSlice = append(strSlice, strings.Trim(strings.Join(columns, ","), ","))
			strSlice = append(strSlice, ") = (")
		} else {
			strSlice = append(strSlice, columns[0])
			strSlice = append(strSlice, " = ")
		}
		valuesToDb := make([]string, len(columns))
		var i int
		for _, val := range columns {
			if _, ok := fields[val]; ok {
				value := fields[val]
				if strings.Contains(value, "'") {
					value = strings.Replace(value, "'", "''", -1)
				}
				valuesToDb[i] = utils.ConcatSlice([]string{"'", value, "'"})
			}
			i++
		}
		if len(fields) > 1 {
			strSlice = append(strSlice, strings.Trim(strings.Join(valuesToDb, ","), ","))
			strSlice = append(strSlice, ") ")
		} else {
			strSlice = append(strSlice, valuesToDb[0])
			strSlice = append(strSlice, " ")
		}
		strSlice = append(strSlice, "WHERE id = ")
		strSlice = append(strSlice, id)
		strSlice = append(strSlice, " RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		var id string
		err := db.QueryRow(queryStr).Scan(&id)
		if err != nil {
			customLog.Logging(err)
		} else {
			response = map[string]string{"id": id}
		}
	}
	return response
}

func (model *Model) GetList(params map[string]string) []map[string]interface{} {
	var resp []map[string]interface{}
	db := customDb.GetConnect()
	defer customDb.CloseConnect(db)
	queryStr := utils.ConcatSlice([]string{
		"SELECT * FROM ",
		model.Table(),
	})
	if filterBy, ok := params["filterBy"]; ok && filterBy != "" {
		queryStr = utils.ConcatSlice([]string{
			queryStr,
			" WHERE ",
			params["filterBy"],
			" = '",
			params["filterVal"],
			"'",
		})
	}
	if order, ok := params["order"]; ok && order != "" {
		queryStr = utils.ConcatSlice([]string{
			queryStr,
			" ORDER BY ",
			params["orderBy"],
			" ",
			params["order"],
		})
	}
	if limit, ok := params["limit"]; ok && limit != "" {
		queryStr = utils.ConcatSlice([]string{
			queryStr,
			" LIMIT ",
			params["limit"],
		})
	}
	if offset, ok := params["offset"]; ok && offset != "" {
		queryStr = utils.ConcatSlice([]string{
			queryStr,
			" OFFSET ",
			params["offset"],
		})
	}
	queryStr = utils.ConcatSlice([]string{
		queryStr,
		" ;",
	})
	rows, err := db.Query(queryStr)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = utils.SqlToMap(rows)
	}
	return resp
}

func (model *Model) GetOneById(id int) map[string]interface{} {
	resp := map[string]interface{}{"success": false, "error": "not found"}
	if id > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"SELECT * FROM ",
			model.Table(),
			" WHERE id=$1;",
		})
		rows, err := db.Query(queryStr, id)
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				resp = data[0]
			}
		}
	}
	return resp
}

func (model *Model) GetOneByTitle(requestData map[string]string) map[string]interface{} {
	resp := map[string]interface{}{"success": false, "error": "not found"}
	if title, ok := requestData["title"]; ok && title != "" {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"SELECT * FROM ",
			model.Table(),
			" WHERE title=$1;",
		})
		rows, err := db.Query(queryStr, requestData["title"])
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				resp = data[0]
			}
		}
	}
	return resp
}

func (model *Model) CheckInterface(v interface{}) bool {
	var resp bool
	if _, ok := v.(HasEvent); ok {
		resp = true
	}
	return resp
}

func (model *Model) Delete(id int) map[string]interface{} {
	resp := map[string]interface{}{"success": false, "error": "not found"}
	if id > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"DELETE FROM ",
			model.Table(),
			" WHERE id=$1 RETURNING id;",
		})
		rows, err := db.Query(queryStr, id)
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				resp = data[0]
			}
		}
	}
	return resp
}

func (model *Model) GetOneCouplet(id, couplet_number int) map[string]interface{} {
	resp := map[string]interface{}{"success": false, "error": "not found"}
	if id > 0 && couplet_number > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"SELECT text FROM ",
			model.Table(),
			" WHERE id=$1 AND text IS NOT NULL;",
		})
		rows, err := db.Query(queryStr, id)
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				if value, ok := data[0]["text"]; ok && value != "" {
					if text, ok := value.(string); ok {
						couplets := strings.Split(string(text), "\n")
						total := len(couplets)
						if couplet_number-1 < total {
							resp = map[string]interface{}{
								"couplet_number": couplet_number,
								"couplet":        couplets[couplet_number-1],
								"total":          total,
							}
						}
					}
				}
			}
		}
	}
	return resp
}
