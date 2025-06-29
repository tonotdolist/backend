package migrate

import "gorm.io/gorm"

var dbMigrationList []any

func RegisterMigrationModel(typ any) {
	dbMigrationList = append(dbMigrationList, typ)
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(dbMigrationList...)
}
