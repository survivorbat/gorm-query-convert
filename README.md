# ⚙️ Gorm Query Convert

[![Go package](https://github.com/survivorbat/gorm-query-convert/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/gorm-query-convert/actions/workflows/test.yaml)

I wanted to provide a map to a WHERE query and automatically convert queries to different operators if characters were present.

By default, all queries are converted, if you want it to be more specific use:

- `SettingOnly()`: Will only change queries on `*gorm.DB` objects that have `.Set("gormqonvert", true)` set.

If you want a particular query to not be converted, use `.Set("gormqonvert", false)`. This works
regardless of configuration.

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
	config := gormqonvert.CharacterConfig{}
	db.Use(gormqonvert.New(config))
}

```

## 🔭 Plans

Not much here.
