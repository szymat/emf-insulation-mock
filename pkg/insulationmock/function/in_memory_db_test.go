package insulationmock_function_test

import (
	"net/http/httptest"
	"testing"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock/function"
	"github.com/yuin/gopher-lua"
)

func TestInMemoryDbSet_FunctionName(t *testing.T) {
	var set insulationmock_function.InMemoryDbSet
	want := "memory_db_set"
	if got := set.FunctionName(); got != want {
		t.Errorf("InMemoryDbSet.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbSetAndGet(t *testing.T) {
	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Set up the mock response writer and request
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// Create and execute InMemoryDbSet
	set := &insulationmock_function.InMemoryDbSet{}
	L.Push(lua.LString("testKey"))
	L.Push(lua.LString("testValue"))
	set.Execute(w, r, L)

	// Verify if the value is set correctly in the memoryDb
	if val, ok := insulationmock_function.MemoryDb["testKey"]; !ok || val != "testValue" {
		t.Errorf("memoryDb did not have the expected value for 'testKey', got %v, want 'testValue'", val)
	}

	// Test InMemoryDbGet
	get := &insulationmock_function.InMemoryDbGet{}
	L.Push(lua.LString("testKey"))
	n := get.Execute(w, r, L)

	// Verify the correct value is pushed onto the Lua stack
	if n != 1 {
		t.Errorf("InMemoryDbGet.Execute() should have returned 1, got %d", n)
	}
	retVal := L.ToString(-1)
	if retVal != "testValue" {
		t.Errorf("Expected to get 'testValue', got '%s'", retVal)
	}
}

func TestInMemoryDbGet_FunctionName(t *testing.T) {
	var get insulationmock_function.InMemoryDbGet
	want := "memory_db_get"
	if got := get.FunctionName(); got != want {
		t.Errorf("InMemoryDbGet.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbDel(t *testing.T) {
	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Set up for the test
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	insulationmock_function.MemoryDb["deleteKey"] = "deleteValue"

	// Test InMemoryDbDel
	del := &insulationmock_function.InMemoryDbDel{}
	L.Push(lua.LString("deleteKey"))
	del.Execute(w, r, L)

	// Verify the key is deleted from memoryDb
	if _, ok := insulationmock_function.MemoryDb["deleteKey"]; ok {
		t.Errorf("memoryDb should not have 'deleteKey' after deletion")
	}
}

func TestInMemoryDbDel_FunctionName(t *testing.T) {
	var del insulationmock_function.InMemoryDbDel
	want := "memory_db_del"
	if got := del.FunctionName(); got != want {
		t.Errorf("InMemoryDbDel.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbFlush(t *testing.T) {
	// Pre-populate memoryDb
	insulationmock_function.MemoryDb["key1"] = "value1"
	insulationmock_function.MemoryDb["key2"] = "value2"

	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Execute InMemoryDbFlush
	flush := &insulationmock_function.InMemoryDbFlush{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	flush.Execute(w, r, L)

	// Verify memoryDb is empty
	if len(insulationmock_function.MemoryDb) != 0 {
		t.Errorf("Expected memoryDb to be empty after flush, found %d items", len(insulationmock_function.MemoryDb))
	}
}

func TestInMemoryDbFlush_FunctionName(t *testing.T) {
	var flush insulationmock_function.InMemoryDbFlush
	want := "memory_db_flush"
	if got := flush.FunctionName(); got != want {
		t.Errorf("InMemoryDbFlush.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbKeys(t *testing.T) {
	// Reset and pre-populate memoryDb
	insulationmock_function.MemoryDb = map[string]string{"key1": "value1", "key2": "value2"}

	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Execute InMemoryDbKeys
	keysFunc := &insulationmock_function.InMemoryDbKeys{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	n := keysFunc.Execute(w, r, L)

	// Verify the correct number of keys is pushed onto the Lua stack
	if n != 1 {
		t.Fatalf("InMemoryDbKeys.Execute() should have returned 1, got %d", n)
	}
	keysTable := L.ToTable(-1)
	if keysTable == nil || keysTable.Len() != 2 {
		t.Errorf("Expected to find 2 keys, found %d", keysTable.Len())
	}
}

func TestInMemoryDbKeys_FunctionName(t *testing.T) {
	var keys insulationmock_function.InMemoryDbKeys
	want := "memory_db_keys"
	if got := keys.FunctionName(); got != want {
		t.Errorf("InMemoryDbKeys.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbValues(t *testing.T) {
	// Reset and pre-populate memoryDb
	insulationmock_function.MemoryDb = map[string]string{"key1": "value1", "key2": "value2"}

	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Execute InMemoryDbValues
	valuesFunc := &insulationmock_function.InMemoryDbValues{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	n := valuesFunc.Execute(w, r, L)

	// Verify the correct number of values is pushed onto the Lua stack
	if n != 1 {
		t.Fatalf("InMemoryDbValues.Execute() should have returned 1, got %d", n)
	}
	valuesTable := L.ToTable(-1)
	if valuesTable == nil || valuesTable.Len() != 2 {
		t.Errorf("Expected to find 2 values, found %d", valuesTable.Len())
	}
}

func TestInMemoryDbValues_FunctionName(t *testing.T) {
	var values insulationmock_function.InMemoryDbValues
	want := "memory_db_values"
	if got := values.FunctionName(); got != want {
		t.Errorf("InMemoryDbValues.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbPairs(t *testing.T) {
	// Reset and pre-populate memoryDb
	insulationmock_function.MemoryDb = map[string]string{"key1": "value1", "key2": "value2"}

	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Execute InMemoryDbPairs
	pairsFunc := &insulationmock_function.InMemoryDbPairs{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	n := pairsFunc.Execute(w, r, L)

	// Verify the correct pairs table is pushed onto the Lua stack
	if n != 1 {
		t.Fatalf("InMemoryDbPairs.Execute() should have returned 1, got %d", n)
	}
	pairsTable := L.ToTable(-1)
	if pairsTable == nil {
		t.Fatalf("Expected to find a table, found nil")
	}
}

func TestInMemoryDbPairs_FunctionName(t *testing.T) {
	var pairs insulationmock_function.InMemoryDbPairs
	want := "memory_db_pairs"
	if got := pairs.FunctionName(); got != want {
		t.Errorf("InMemoryDbPairs.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbHas_FunctionName(t *testing.T) {
	var has insulationmock_function.InMemoryDbHas
	want := "memory_db_has"
	if got := has.FunctionName(); got != want {
		t.Errorf("InMemoryDbHas.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbHas(t *testing.T) {
	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Set up for the test
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	insulationmock_function.MemoryDb["testKey"] = "testValue"

	// Test InMemoryDbHas
	has := &insulationmock_function.InMemoryDbHas{}
	L.Push(lua.LString("testKey"))
	n := has.Execute(w, r, L)

	// Verify the correct value is pushed onto the Lua stack
	if n != 1 {
		t.Errorf("InMemoryDbHas.Execute() should have returned 1, got %d", n)
	}
	retVal := L.ToBool(-1)
	if !retVal {
		t.Errorf("Expected to get true, got false")
	}
}

func TestInMemoryDbSize_FunctionName(t *testing.T) {
	var size insulationmock_function.InMemoryDbSize
	want := "memory_db_size"
	if got := size.FunctionName(); got != want {
		t.Errorf("InMemoryDbSize.FunctionName() = %v, want %v", got, want)
	}
}

func TestInMemoryDbSize(t *testing.T) {
	// Initialize the Lua state
	L := lua.NewState()
	defer L.Close()

	// Set up for the test
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	insulationmock_function.MemoryDb = make(map[string]string)
	insulationmock_function.MemoryDb["key1"] = "value1"
	insulationmock_function.MemoryDb["key2"] = "value2"

	// Test InMemoryDbSize
	size := &insulationmock_function.InMemoryDbSize{}
	n := size.Execute(w, r, L)

	// Verify the correct size is pushed onto the Lua stack
	if n != 1 {
		t.Fatalf("InMemoryDbSize.Execute() should have returned 1, got %d", n)
	}
	retVal := L.ToNumber(-1)
	if retVal != 2 {
		t.Errorf("Expected to get 2, got %d", int(retVal))
	}
}
