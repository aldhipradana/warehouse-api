// restful/interfaces.go
package restful

// Filterable ensures the model tells the controller which columns are searchable.
// Equivalent to getSearchable() in your PHP Trait.
type Filterable interface {
	GetSearchableFields() []string
}

// FilterRequest represents the JSON structure of your specific filters
// e.g. ?filter={"status": "active", "created_at": {"operator": ">", "value": "..."}}
type FilterValue struct {
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Function string      `json:"function"` // date, time, in, between
}
