package gormqonvert

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const tagName = "gormQonvert"

func (d *gormQonvert) queryCallback(db *gorm.DB) {
	// If we only want to like queries that are explicitly set to true, we back out early if anything's amiss
	settingValue, settingOk := db.Get(tagName)
	if d.conditionalSetting && !settingOk {
		return
	}

	if settingOk {
		if boolValue, _ := settingValue.(bool); !boolValue {
			return
		}
	}

	exp, settingOk := db.Statement.Clauses["WHERE"].Expression.(clause.Where)
	if !settingOk {
		return
	}

	for index, cond := range exp.Exprs {
		switch cond := cond.(type) {
		case clause.Eq:
			value, ok := cond.Value.(string)
			if !ok {
				continue
			}

			var condition string

			switch {
			case d.config.GreaterOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.GreaterOrEqualToPrefix):
				condition = fmt.Sprintf("%s >= ?", cond.Column)
				value = value[len(d.config.GreaterOrEqualToPrefix):]

			case d.config.GreaterThanPrefix != "" && strings.HasPrefix(value, d.config.GreaterThanPrefix):
				condition = fmt.Sprintf("%s > ?", cond.Column)
				value = value[len(d.config.GreaterThanPrefix):]

			case d.config.LessOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.LessOrEqualToPrefix):
				condition = fmt.Sprintf("%s <= ?", cond.Column)
				value = value[len(d.config.LessOrEqualToPrefix):]

			case d.config.LessThanPrefix != "" && strings.HasPrefix(value, d.config.LessThanPrefix):
				condition = fmt.Sprintf("%s < ?", cond.Column)
				value = value[len(d.config.LessThanPrefix):]

			case d.config.NotEqualToPrefix != "" && strings.HasPrefix(value, d.config.NotEqualToPrefix):
				condition = fmt.Sprintf("%s != ?", cond.Column)
				value = value[len(d.config.NotEqualToPrefix):]

			case d.config.LikePrefix != "" && strings.HasPrefix(value, d.config.LikePrefix):
				condition = fmt.Sprintf("%s LIKE ?", cond.Column)
				value = value[len(d.config.LikePrefix):]

			case d.config.NotLikePrefix != "" && strings.HasPrefix(value, d.config.NotLikePrefix):
				condition = fmt.Sprintf("%s NOT LIKE ?", cond.Column)
				value = value[len(d.config.NotLikePrefix):]

			default:
				continue
			}

			exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(condition, value).Statement.Clauses["WHERE"].Expression
		case clause.IN:
			var conversionCounter int
			var useOr bool

			query := db.Session(&gorm.Session{NewDB: true})

			for _, value := range cond.Values {
				value, ok := value.(string)
				if !ok {
					continue
				}

				condition := fmt.Sprintf("%s = ?", cond.Column)

				switch {
				case d.config.GreaterOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.GreaterOrEqualToPrefix):
					condition = fmt.Sprintf("%s >= ?", cond.Column)
					value = value[len(d.config.GreaterOrEqualToPrefix):]

				case d.config.GreaterThanPrefix != "" && strings.HasPrefix(value, d.config.GreaterThanPrefix):
					condition = fmt.Sprintf("%s > ?", cond.Column)
					value = value[len(d.config.GreaterThanPrefix):]

				case d.config.LessOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.LessOrEqualToPrefix):
					condition = fmt.Sprintf("%s <= ?", cond.Column)
					value = value[len(d.config.LessOrEqualToPrefix):]

				case d.config.LessThanPrefix != "" && strings.HasPrefix(value, d.config.LessThanPrefix):
					condition = fmt.Sprintf("%s < ?", cond.Column)
					value = value[len(d.config.LessThanPrefix):]

				case d.config.NotEqualToPrefix != "" && strings.HasPrefix(value, d.config.NotEqualToPrefix):
					condition = fmt.Sprintf("%s != ?", cond.Column)
					value = value[len(d.config.NotEqualToPrefix):]

				case d.config.LikePrefix != "" && strings.HasPrefix(value, d.config.LikePrefix):
					condition = fmt.Sprintf("%s LIKE ?", cond.Column)
					value = value[len(d.config.LikePrefix):]

				case d.config.NotLikePrefix != "" && strings.HasPrefix(value, d.config.NotLikePrefix):
					condition = fmt.Sprintf("%s NOT LIKE ?", cond.Column)
					value = value[len(d.config.NotLikePrefix):]

				default:
					continue
				}

				conversionCounter++
				if useOr {
					query = query.Or(condition, value)
					continue
				}

				query = query.Where(condition, value)
				useOr = true
			}

			// Don't alter the query if it isn't necessary
			if conversionCounter == 0 {
				continue
			}

			exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(query).Statement.Clauses["WHERE"].Expression
		}
	}
}
