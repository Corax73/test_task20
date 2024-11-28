package repository

import (
	"songLibrary/customDb"
	"songLibrary/models"
)

type Repository struct {

}

// Init performs initial migrations to validate and create tables. Returns true on success.
func (repository *Repository) Init(modelsList []*models.Model) bool {
	var resp bool
	if len(modelsList) > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		for _, model := range modelsList {
			if resp = model.CheckModelTable(db); !resp {
				if resp = model.RunTableMigration(db); !resp {
					break
				}
			}
		}
	}
	return resp
}