package insulationmock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock/function"
	"github.com/yuin/gopher-lua"
)

// LuaRouter handles routing of requests to Lua scripts
type LuaRouter struct {
	L      *lua.LState
	Config *config
	routes []Route
}

// NewLuaRouter creates and initializes a new LuaRouter
func NewLuaRouter(config *config) *LuaRouter {
	L := lua.NewState()
	defer L.Close()
	return &LuaRouter{
		L:      L,
		Config: config,
	}
}

// LoadRoutes loads and parses route definitions from Lua scripts
func (lr *LuaRouter) LoadRoutes() error {
	// Open the Lua script file
	lr.L = lua.NewState()

	defer lr.L.Close()

	// Load the script
	if err := lr.L.DoFile(lr.Config.ScriptsPath + lr.Config.RoutesPath); err != nil {
		log.Fatal(err)
	}

	// Example of parsing routes from Lua, adjust as needed
	routesTable := lr.L.GetGlobal("Routers").(*lua.LTable)
	routesTable.ForEach(func(_, value lua.LValue) {
		route := value.(*lua.LTable)

		if route.RawGetString("pattern").Type() == lua.LTNil {
			log.Fatal("Invalid route configuration")
		}
		pattern := route.RawGetString("pattern").String()
		// check if route pattern contains * or ? and replace it with regex
		pattern = regexp.MustCompile(`\*`).ReplaceAllString(pattern, `.*`)
		pattern = regexp.MustCompile(`\?`).ReplaceAllString(pattern, `.`)

		if route.RawGetString("methods").Type() == lua.LTNil {
			log.Fatal("Invalid route configuration")
		}

		methodsTable := route.RawGetString("methods").(*lua.LTable)
		var methods []string
		methodsTable.ForEach(func(_, methodValue lua.LValue) {
			methods = append(methods, methodValue.String())
		})
		priority := int(route.RawGetString("priority").(lua.LNumber))
		if script := route.RawGetString("script"); script.Type() != lua.LTNil {
			lr.routes = append(lr.routes, routeEntryScripted{
				routeEntry: routeEntry{
					pattern:  pattern,
					priority: priority,
					methods:  methods,
				},
				script: script.String(),
			})
		} else if response := route.RawGetString("response"); response.Type() == lua.LTTable {
			statusCode := response.(*lua.LTable).RawGetInt(1).String()
			responseBody := response.(*lua.LTable).RawGetInt(2).String()

			statusCodeInt, _ := strconv.Atoi(statusCode)

			lr.routes = append(lr.routes, routeEntryResponse{
				routeEntry: routeEntry{
					pattern:  pattern,
					priority: priority,
					methods:  methods,
				},
				response: responseType{
					statusCode: statusCodeInt,
					body:       responseBody,
				},
			})
		} else {
			log.Fatal("Invalid route configuration")
		}
	})

	return nil // Or relevant error
}

func (lr *LuaRouter) ExecuteScriptFile(w http.ResponseWriter, r *http.Request, scriptPath string) {
	if _, err := os.Stat(lr.Config.ScriptsPath + "/" + scriptPath); os.IsNotExist(err) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	scriptContent, err := os.ReadFile(lr.Config.ScriptsPath + "/" + scriptPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	lr.ExecuteScript(w, r, string(scriptContent))
}

// ExecuteScript executes the Lua script for the given route and handles the HTTP response.
func (lr *LuaRouter) ExecuteScript(w http.ResponseWriter, r *http.Request, script string) {
	// Check if the script exists
	l := lua.NewState()
	defer l.Close()

	// Set lua global "method" and "uri", with "headers"
	l.SetGlobal("method", lua.LString(r.Method))
	l.SetGlobal("uri", lua.LString(r.URL.Path))
	headers := l.NewTable()
	for key, val := range r.Header {
		headers.RawSetString(key, lua.LString(val[0]))
	}
	l.SetGlobal("headers", headers)

	responseHeaders := map[string]string{}

	luaFunctions := []insulationmock_function.LuaFunction{
		&insulationmock_function.AddHeader{},
		&insulationmock_function.InMemoryDbSet{},
		&insulationmock_function.InMemoryDbGet{},
		&insulationmock_function.InMemoryDbDel{},
		&insulationmock_function.InMemoryDbFlush{},
		&insulationmock_function.InMemoryDbKeys{},
		&insulationmock_function.InMemoryDbValues{},
		&insulationmock_function.InMemoryDbPairs{},
	}
	// Add all function to lua that implements the lua_function interface
	for _, function := range luaFunctions {
		l.SetGlobal(function.FunctionName(), l.NewFunction(func(L *lua.LState) int {
			return function.Execute(w, r, L)
		}))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// pass the body/content to lua scripts
	l.SetGlobal("body", lua.LString(string(body)))

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err == nil {
	
	luaData := l.NewTable()
	for key, value := range data {
		l.SetField(luaData, key, lua.LString(fmt.Sprintf("%v", value)))
	}
	l.SetGlobal("decoded_body", luaData)
	} else {
		l.SetGlobal("decoded_body", l.NewTable())
	}

	if _, err := os.Stat(lr.Config.BeforeScript); !os.IsNotExist(err) {
		if l.DoFile(lr.Config.BeforeScript); err != nil {
			log.Printf("Error executing Lua script: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err := l.DoString(script); err != nil {
		// if err := lr.L.DoFile(lr.Config.ScriptsPath + "/" + scriptPath); err != nil {
		log.Printf("Error executing Lua script: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// check if after script exists
	if _, err := os.Stat(lr.Config.AfterScript); !os.IsNotExist(err) {
		if l.DoFile(lr.Config.AfterScript); err != nil {
			log.Printf("Error executing Lua script: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Get the response
	status := l.ToInt(-2)
	resp := l.Get(-1)
	l.Pop(2)

	// Write the response

	for key, value := range responseHeaders {
		w.Header().Set(key, value)
	}
	w.WriteHeader(status)
	w.Write([]byte(resp.String()))
}

// SortRoutes sorts the loaded routes by their priority
func (lr *LuaRouter) SortRoutes() {
	sort.Slice(lr.routes, func(i, j int) bool {
		return lr.routes[i].GetPriority() < lr.routes[j].GetPriority()
	})
}

func (lr *LuaRouter) FindAndExecuteRoute(w http.ResponseWriter, r *http.Request) {
	for _, route := range lr.routes {
		if route.MatchUrl(r.URL.Path, r.Method) {
			log.Printf("%s %s [%d|%s] matched", r.Method, r.URL.Path, route.GetPriority(), route.GetPattern())
			if responseRoute, ok := route.(routeEntryResponse); ok {
				composed := fmt.Sprintf("return %d, \"%s\"", responseRoute.response.statusCode, responseRoute.response.body)
				lr.ExecuteScript(w, r, composed)
				return
			}

			if scriptedRoute, ok := route.(routeEntryScripted); ok {
				lr.ExecuteScriptFile(w, r, scriptedRoute.script)
				return
			}

		} else {
			log.Printf("%s %s [%d|%s] not matched", r.Method, r.URL.Path, route.GetPriority(), route.GetPattern())
		}
	}

	// No matching route found; return 404 Not Found.
	http.NotFound(w, r)
}
