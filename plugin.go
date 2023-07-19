package gormqonvert

import (
	"gorm.io/gorm"
)

// Compile-time interface check
var _ gorm.Plugin = new(gormQonvert)

// Option can be given to the New() method to tweak its behaviour
type Option func(like *gormQonvert)

// CharacterConfig must be provided to indicate when a field must be changed, these should all be prefixes
type CharacterConfig struct {
	GreaterThanPrefix      string
	GreaterOrEqualToPrefix string
	LessThanPrefix         string
	LessOrEqualToPrefix    string
	NotEqualToPrefix       string
}

// SettingOnly makes it so that only queries with the setting 'gormQonvert' set to true can be turned into LIKE queries.
// This can be configured using db.Set("gormQonvert", true) on the query.
func SettingOnly() Option {
	return func(like *gormQonvert) {
		like.conditionalSetting = true
	}
}

// New creates a new instance of the plugin that can be registered in gorm. Without any settings, all queries will be
// LIKE-d.
func New(config CharacterConfig, opts ...Option) gorm.Plugin {
	plugin := &gormQonvert{config: config}

	for _, opt := range opts {
		opt(plugin)
	}

	return plugin
}

type gormQonvert struct {
	conditionalSetting bool

	config CharacterConfig
}

func (d *gormQonvert) Name() string {
	return "gormQonvert"
}

func (d *gormQonvert) Initialize(db *gorm.DB) error {
	return db.Callback().Query().Before("gorm:query").Register("gormQonvert:query", d.queryCallback)
}
