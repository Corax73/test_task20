package customDb

import (
	"songLibrary/models"
)

// Init performs initial migrations to validate and create tables. Returns true on success.
func Init(modelsList []*models.Model) bool {
	var resp bool
	if len(modelsList) > 0 {
		db := GetConnect()
		defer CloseConnect(db)
		for _, model := range modelsList {
			if !model.CheckModelTable(db) {
				if resp = model.RunTableMigration(db); !resp {
					break
				}
			}
		}
	}
	return resp
}
