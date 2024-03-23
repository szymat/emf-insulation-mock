package insulationmock_function_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock/function"
	"github.com/yuin/gopher-lua"
)

func TestAddHeader_FunctionName(t *testing.T) {
	var a insulationmock_function.AddHeader
	want := "add_header"
	if got := a.FunctionName(); got != want {
		t.Errorf("AddHeader.FunctionName() = %v, want %v", got, want)
	}
}

type MockResponseWriter struct {
	header http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{header: http.Header{}}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.header
}

func (m *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {}

func TestAddHeader_Execute(t *testing.T) {
	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Set up the mock response writer and request
	w := NewMockResponseWriter()
	r := httptest.NewRequest("GET", "/", bytes.NewBufferString(""))

	// Define the key and value to be set in the header
	key := "Content-Type"
	value := "application/json"

	// Push key and value onto the Lua stack
	L.Push(lua.LString(key))
	L.Push(lua.LString(value))

	// Execute the AddHeader function
	a := insulationmock_function.AddHeader{}
	a.Execute(w, r, L)

	// Check if the header was set correctly
	if got := w.Header().Get(key); got != value {
		t.Errorf("Execute() did not set the correct header, got %v, want %v", got, value)
	}
}
