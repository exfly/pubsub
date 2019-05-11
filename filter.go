package pubsub

type Filter struct {
	filter map[string]map[string]bool
}

func NewFilter() *Filter {
	return &Filter{filter: map[string]map[string]bool{}}
}
func (ef *Filter) Len() int {
	return len(ef.filter)
}
func (ef *Filter) Add(key, val string) {
	if _, ok := ef.filter[key]; ok {
		ef.filter[key][val] = true
	} else {
		ef.filter[key] = map[string]bool{val: true}
	}
	return
}

// Include event filter
func (ef *Filter) Include(msg Msg) bool {
	return ef.matchScope(msg) &&
		ef.matchAction(msg)
}
func (ef *Filter) matchScope(msg Msg) bool {
	if !ef.Contains("scope") {
		return true
	}
	return ef.ExactMatch("scope", msg.Scope)
}
func (ef *Filter) matchAction(msg Msg) bool {
	if !ef.Contains("action") {
		return true
	}
	return ef.ExactMatch("action", msg.Action)
}

// Contains returns true if the key exists in the mapping
func (ef *Filter) Contains(key string) bool {
	_, ok := ef.filter[key]
	return ok
}

// ExactMatch returns true if the source matches exactly one of the values.
func (ef *Filter) ExactMatch(key, source string) bool {
	fieldValues, ok := ef.filter[key]
	//do not filter if there is no filter set or cannot determine filter
	if !ok || len(fieldValues) == 0 {
		return true
	}
	// try to match full name value to avoid O(N) regular expression matching
	return fieldValues[source]
}
