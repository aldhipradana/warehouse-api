// restful/scopes.go
package restful

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ApplyFilters corresponds to your protected function filterAll()
func ApplyFilters(db *gorm.DB, filterJSON string, search string, model interface{}) *gorm.DB {
	// 1. Handle Global Search (q)
	if search != "" {
		if m, ok := model.(Filterable); ok {
			fields := m.GetSearchableFields()
			if len(fields) > 0 {
				db = db.Where(buildSearchQuery(fields, search))
			}
		}
	}

	// 2. Handle JSON Filters
	if filterJSON == "" {
		return db
	}

	var filters map[string]interface{}
	if err := json.Unmarshal([]byte(filterJSON), &filters); err != nil {
		return db // Ignore invalid JSON
	}

	for key, rawVal := range filters {
		// Handle Relations (e.g., "schedules.schedule_id")
		if strings.Contains(key, ".") {
			parts := strings.Split(key, ".")
			relation := parts[0]
			col := parts[1]

			// This replicates: whereHas('relation', function($q) { ... })
			// Note: GORM uses Joins or association queries differently.
			// The simplest GORM equivalent for filtering on relations:
			db = db.Joins(relation).Where(fmt.Sprintf("%s.%s = ?", relation, col), rawVal)
			continue
		}

		// Parse the value (it could be a raw string or a JSON object)
		valMap, isObj := rawVal.(map[string]interface{})

		if !isObj {
			// Simple equality: "status": "active"
			db = db.Where(fmt.Sprintf("%s = ?", key), rawVal)
			continue
		}

		// Complex logic: operator, function, between, etc.
		operator := "="
		if op, ok := valMap["operator"].(string); ok {
			operator = op
		}
		value := valMap["value"]
		fn := ""
		if f, ok := valMap["function"].(string); ok {
			fn = f
		}

		// Handle Functions (date, in, between)
		switch fn {
		case "date":
			db = db.Where(fmt.Sprintf("DATE(%s) %s ?", key, operator), value)
		case "in":
			// value should be "1,2,3"
			strVal := fmt.Sprintf("%v", value)
			db = db.Where(fmt.Sprintf("%s IN ?", key), strings.Split(strVal, ","))
		case "between":
			strVal := fmt.Sprintf("%v", value)
			rangeVals := strings.Split(strVal, ",")
			if len(rangeVals) == 2 {
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", key), rangeVals[0], rangeVals[1])
			}
		case "like":
			db = db.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%v%%", value))
		default:
			// Standard Operator
			db = db.Where(fmt.Sprintf("%s %s ?", key, operator), value)
		}
	}

	return db
}

func buildSearchQuery(fields []string, search string) string {
	// Builds: field1 LIKE %s% OR field2 LIKE %s%
	var query []string
	for _, field := range fields {
		query = append(query, fmt.Sprintf("%s LIKE '%%%s%%'", field, search))
	}
	return strings.Join(query, " OR ")
}
