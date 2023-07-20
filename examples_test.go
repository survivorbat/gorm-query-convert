package gormqonvert

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ExampleNew() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	config := CharacterConfig{
		GreaterThanPrefix:      ">",
		GreaterOrEqualToPrefix: ">=",
		LessThanPrefix:         "<",
		LessOrEqualToPrefix:    "<=",
		NotEqualToPrefix:       "!=",
		LikePrefix:             "~",
		NotLikePrefix:          "!~",
	}

	_ = db.Use(New(config))

	_ = db.Use(New(config, SettingOnly()))
}
