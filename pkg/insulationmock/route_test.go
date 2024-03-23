package insulationmock_test

import (
	"testing"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock"
)

func TestRouteMatching(t *testing.T) {
	route := insulationmock.NewRoute("/test", 1, []string{"GET", "POST"})

	tests := []struct {
		path   string
		method string
		match  bool
	}{
		{"/test", "GET", true},
		{"/test", "POST", true},
		{"/test", "PUT", false},
		{"/doesnotmatch", "GET", false},
	}

	for _, tc := range tests {
		if got := route.MatchUrl(tc.path, tc.method); got != tc.match {
			t.Errorf("MatchUrl(%q, %q) = %v; want %v", tc.path, tc.method, got, tc.match)
		}
	}
}

func TestRoutePriority(t *testing.T) {
	priority := 5
	route := insulationmock.NewRoute("/test", priority, []string{"GET"})
	if got := route.GetPriority(); got != priority {
		t.Errorf("GetPriority() = %d; want %d", got, priority)
	}
}

func TestRoutePattern(t *testing.T) {
	pattern := "/test"
	route := insulationmock.NewRoute(pattern, 1, []string{"GET"})
	if got := route.GetPattern(); got != pattern {
		t.Errorf("GetPattern() = %s; want %s", got, pattern)
	}
}
