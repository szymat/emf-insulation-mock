package insulationmock_function

import (
	"net/http"

	"github.com/yuin/gopher-lua"
)

type AddHeader struct{}

func (a *AddHeader) FunctionName() string {
	return "add_header"
}

func (a *AddHeader) Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int {
	key := L.ToString(1)
	value := L.ToString(2)
	w.Header().Set(key, value)
	return 0
}
