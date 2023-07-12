package persistence

type QueryFilter struct {
	Filters map[string]string
	Sort    map[string]string
}
