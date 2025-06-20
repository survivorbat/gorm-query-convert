package gormqonvert

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const tagName = "gormQonvert"

func (d *gormQonvert) replaceExpressions(db *gorm.DB, expressions []clause.Expression) []clause.Expression {
	for index, cond := range expressions {
		switch cond := cond.(type) {
		case clause.AndConditions:
			// Recursively go through the expressions of AndConditions
			cond.Exprs = d.replaceExpressions(db, cond.Exprs)
			expressions[index] = cond
		case clause.OrConditions:
			// Recursively go through the expressions of OrConditions
			cond.Exprs = d.replaceExpressions(db, cond.Exprs)
			expressions[index] = cond
		case clause.Eq:
			value, ok := cond.Value.(string)
			if !ok {
				continue
			}

			var condition string

			column := cond.Column.(clause.Column)

			switch {
			case d.config.GreaterOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.GreaterOrEqualToPrefix):
				condition = fmt.Sprintf("%s >= ?", column.Name)
				value = value[len(d.config.GreaterOrEqualToPrefix):]

			case d.config.GreaterThanPrefix != "" && strings.HasPrefix(value, d.config.GreaterThanPrefix):
				condition = fmt.Sprintf("%s > ?", column.Name)
				value = value[len(d.config.GreaterThanPrefix):]

			case d.config.LessOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.LessOrEqualToPrefix):
				condition = fmt.Sprintf("%s <= ?", column.Name)
				value = value[len(d.config.LessOrEqualToPrefix):]

			case d.config.LessThanPrefix != "" && strings.HasPrefix(value, d.config.LessThanPrefix):
				condition = fmt.Sprintf("%s < ?", column.Name)
				value = value[len(d.config.LessThanPrefix):]

			case d.config.NotEqualToPrefix != "" && strings.HasPrefix(value, d.config.NotEqualToPrefix):
				condition = fmt.Sprintf("%s != ?", column.Name)
				value = value[len(d.config.NotEqualToPrefix):]

			case d.config.LikePrefix != "" && strings.HasPrefix(value, d.config.LikePrefix):
				condition = fmt.Sprintf("%s LIKE ?", column.Name)
				value = value[len(d.config.LikePrefix):]

			case d.config.NotLikePrefix != "" && strings.HasPrefix(value, d.config.NotLikePrefix):
				condition = fmt.Sprintf("%s NOT LIKE ?", column.Name)
				value = value[len(d.config.NotLikePrefix):]

			default:
				continue
			}

			expressions[index] = db.Session(&gorm.Session{NewDB: true}).Where(condition, value).Statement.Clauses["WHERE"].Expression
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

				column := cond.Column.(clause.Column)

				switch {
				case d.config.GreaterOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.GreaterOrEqualToPrefix):
					condition = fmt.Sprintf("%s >= ?", column.Name)
					value = value[len(d.config.GreaterOrEqualToPrefix):]

				case d.config.GreaterThanPrefix != "" && strings.HasPrefix(value, d.config.GreaterThanPrefix):
					condition = fmt.Sprintf("%s > ?", column.Name)
					value = value[len(d.config.GreaterThanPrefix):]

				case d.config.LessOrEqualToPrefix != "" && strings.HasPrefix(value, d.config.LessOrEqualToPrefix):
					condition = fmt.Sprintf("%s <= ?", column.Name)
					value = value[len(d.config.LessOrEqualToPrefix):]

				case d.config.LessThanPrefix != "" && strings.HasPrefix(value, d.config.LessThanPrefix):
					condition = fmt.Sprintf("%s < ?", column.Name)
					value = value[len(d.config.LessThanPrefix):]

				case d.config.NotEqualToPrefix != "" && strings.HasPrefix(value, d.config.NotEqualToPrefix):
					condition = fmt.Sprintf("%s != ?", column.Name)
					value = value[len(d.config.NotEqualToPrefix):]

				case d.config.LikePrefix != "" && strings.HasPrefix(value, d.config.LikePrefix):
					condition = fmt.Sprintf("%s LIKE ?", column.Name)
					value = value[len(d.config.LikePrefix):]

				case d.config.NotLikePrefix != "" && strings.HasPrefix(value, d.config.NotLikePrefix):
					condition = fmt.Sprintf("%s NOT LIKE ?", column.Name)
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

			expressions[index] = db.Session(&gorm.Session{NewDB: true}).Where(query).Statement.Clauses["WHERE"].Expression
		}
	}
	return expressions
}

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

	exp.Exprs = d.replaceExpressions(db, exp.Exprs)
}
