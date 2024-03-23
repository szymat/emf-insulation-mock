package insulationmock

import (
	"regexp"
)

// route defines the interface for routing to Lua scripts or static responses.
type Route interface {
	MatchUrl(path string, method string) bool
	GetPriority() int
	GetPattern() string
}


func NewRoute(pattern string, priority int, methods []string) Route {
	return &routeEntry{
		pattern:  pattern,
		priority: priority,
		methods:  methods,
	}
}

// routeEntry struct
type routeEntry struct {
	pattern  string
	priority int
	methods  []string
}

type routeEntryScripted struct {
	routeEntry
	script string
}

type routeEntryResponse struct {
	routeEntry
	response responseType
}

type responseType struct {
	statusCode int
	body       string
}

func contains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func (r routeEntry) MatchUrl(path string, method string) bool {
	return regexp.MustCompile("^"+r.pattern+"$").MatchString(path) && (contains(r.methods, method) || contains(r.methods, "ALL"))
}

func (r routeEntry) GetPriority() int {
	return r.priority
}

func (r routeEntry) GetPattern() string {
	return r.pattern
}
