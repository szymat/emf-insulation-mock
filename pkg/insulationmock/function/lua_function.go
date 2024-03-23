package insulationmock_function

import (
	"net/http"

	"github.com/yuin/gopher-lua"
)

type LuaFunction interface {
	FunctionName() string
	Execute(w http.ResponseWriter, r *http.Request, L *lua.LState) int
}
