# ⚙️ Gorm Query Convert

[![Go package](https://github.com/survivorbat/gorm-query-convert/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/gorm-query-convert/actions/workflows/test.yaml)

Laziness rules, why write GORM queries if you can simply add prefixes to a `map[string]any`'s values and automatically
convert queries to use different operators. All prefix characters can be custom-defined and are only enabled if you define them.
The currently supported queries are:

- `WHERE x != y`
- `WHERE x >= y`
- `WHERE x > y`
- `WHERE x <= y`
- `WHERE x < y`
- `WHERE x LIKE y`
- `WHERE x NOT LIKE y`

By default, all queries are converted, if you want it to be more specific use:

- `SettingOnly()`: Will only change queries on `*gorm.DB` objects that have `.Set("gormqonvert", true)` set.

If you want a particular query to not be converted, use `.Set("gormqonvert", false)`. This works
regardless of configuration.

## 💡 Related Libraries 

- [deepgorm](https://github.com/survivorbat/gorm-deep-filtering) turns nested maps in WHERE-calls into subqueries
- [gormlike](https://github.com/survivorbat/gorm-like) turns WHERE-calls into LIkE queries if certain tokens were found
- [gormcase](https://github.com/survivorbat/gorm-case) adds case insensitivity to WHERE queries
- [gormtestutil](https://github.com/ing-bank/gormtestutil) provides easy utility methods for unit-testing with gorm

## ⬇️ Installation

`go get github.com/survivorbat/gorm-query-convert`

## 📋 Usage

```go
package main

import (
    "github.com/survivorbat/gorm-query-convert"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	config := gormqonvert.CharacterConfig{
		GreaterThanPrefix:      ">",
		GreaterOrEqualToPrefix: ">=",
		LessThanPrefix:         "<",
		LessOrEqualToPrefix:    "<=",
		NotEqualToPrefix:       "!=",
		LikePrefix:             "~",
		NotLikePrefix:          "!~",
    }
	db.Use(gormqonvert.New(config))
}

```

## 🔭 Plans

Not much here.
