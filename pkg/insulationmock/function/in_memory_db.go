package insulationmock_function

import (
	"net/http"

	"github.com/yuin/gopher-lua"
)

var MemoryDb = make(map[string]string)

type InMemoryDbSet struct{}

func (a *InMemoryDbSet) FunctionName() string {
	return "memory_db_set"
}

func (a *InMemoryDbSet) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	key := L.ToString(1)
	value := L.ToString(2)
	MemoryDb[key] = value
	return 0
}

type InMemoryDbGet struct{}

func (a *InMemoryDbGet) FunctionName() string {
	return "memory_db_get"
}

func (a *InMemoryDbGet) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	key := L.ToString(1)
	L.Push(lua.LString(MemoryDb[key]))
	return 1
}

type InMemoryDbDel struct{}

func (a *InMemoryDbDel) FunctionName() string {
	return "memory_db_del"
}

func (a *InMemoryDbDel) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	key := L.ToString(1)
	delete(MemoryDb, key)
	return 0
}

type InMemoryDbFlush struct{}

func (a *InMemoryDbFlush) FunctionName() string {
	return "memory_db_flush"
}

func (a *InMemoryDbFlush) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	MemoryDb = make(map[string]string)
	return 0
}

type InMemoryDbKeys struct{}

func (a *InMemoryDbKeys) FunctionName() string {
	return "memory_db_keys"
}

func (a *InMemoryDbKeys) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	keys := L.NewTable()
	i := 1
	for key := range MemoryDb {
		L.RawSet(keys, lua.LNumber(i), lua.LString(key))
		i++
	}
	L.Push(keys)
	return 1
}

type InMemoryDbValues struct{}

func (a *InMemoryDbValues) FunctionName() string {
	return "memory_db_values"
}

func (a *InMemoryDbValues) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	values := L.NewTable()
	i := 1
	for _, value := range MemoryDb {
		L.RawSet(values, lua.LNumber(i), lua.LString(value))
		i++
	}
	L.Push(values)
	return 1
}

// All pairs

type InMemoryDbPairs struct{}

func (a *InMemoryDbPairs) FunctionName() string {
	return "memory_db_pairs"
}

func (a *InMemoryDbPairs) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	// return as key: value dictionary in lua
	pairs := L.NewTable()
	for key, value := range MemoryDb {
		pair := L.NewTable()
		L.SetField(pair, "key", lua.LString(key))
		L.SetField(pair, "value", lua.LString(value))
		L.RawSet(pairs, lua.LString(key), pair)
	}
	L.Push(pairs)
	return 1
}
