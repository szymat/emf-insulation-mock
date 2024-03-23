# EMF Insulation Mock Server

The EMF Insulation Mock Server is designed to provide a flexible and easy-to-use platform for mocking HTTP services. It uses Lua scripts to define the behavior of the mock server, allowing dynamic responses to incoming requests based on the request path, method, and body content.

## Features

- **Dynamic Routing**: Define routes and responses using Lua scripts for maximum flexibility.
- **Pre and Post Processing**: Execute Lua scripts before and after handling requests to perform additional processing or logging.
- **Simple Configuration**: Easy configuration through command-line flags or environment variables.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Lua 5.3 (if running Lua scripts externally)

### Installation

Clone the repository to your local machine:

```bash
git clone https://github.com/szymat/emf-insulation-mock.git
cd emf-insulation-mock
```

### Build the server:
    
```bash
go build -o emf-insulation-mock ./cmd/volt-emf
```

### Running the Server
To start the server, run:

```bash
./emf-insulation-mock -scripts path/to/your/lua/scripts
```


## Command-Line Flags

- scripts: Path to the directory containing the Lua scripts. Default is `scripts/`.
- routes: Name of the Lua file containing route definitions. Default is `routes.lua`.
- before: Name of the Lua file containing the before script. Default is `before.lua`.
- after: Name of the Lua file containing the after script. Default is `after.lua`.
- port: Port on which the server will listen. Default is `8081`.

## Configuring routes

Routes are defined in scripts/routes.lua. The file should return a table containing route definitions. 
Each route definition is a table with the following fields:

- pattern: A string pattern or regex that matches the request path.
- script or response - one of them is required to define the route. Not both.
- script: Name of the Lua script to execute for this route. `script = "string.lua"` will execute the script `scripts/string.lua`. The script should return a table with the following fields:
- response: A table containing the response status code and body. `response = { 200, "response body" }`

- methods: An array of HTTP methods that this route should match. `methods = { "GET", "POST" }`
- priority: An integer that determines the order in which routes are matched. Lower numbers have higher priority. Routes with the same priority are matched in the order they are defined (or random because of Lua)

### Example routes.lua

`scripts/routes.lua`:
```lua
Routes = {
  { pattern = "string/regex/*", script = "string.lua", response = { 200, "response body" }, methods = { "GET", "POST" }, priority = 1 }
  -- ..
  { pattern = "/services/oauth2/token", script = "oauth.lua",            methods = { "POST", "PUT", "PATCH" },  priority = 1 },
  { pattern = "*",                      response = { 404, "Not found" }, methods = { "GET" },                   priority = 100 },
  { pattern = "*",                      response = { 201, "" },          methods = { "POST", "PUT", "DELETE" }, priority = 101 },
  { pattern = "/",                      response = { 200, "" },          methods = { "ALL" },                   priority = 2 },
  { pattern = "/article/*",             response = { 201, "okey" },      methods = { "ALL" },                   priority = 3 }
}

return Routers
```

`scripts/oauth.lua`:
```lue
if body == nil then
  return 403, "no body"
end

local clientId = decoded_body["clientId"]
local clientSecret = decoded_body["clientSecret"]

-- print table
for k, v in pairs(decoded_body) do
  print(k, v)
end

if clientId == nil or clientSecret == nil then
  return 403, "no clientId or clientSecret"
end

memory_db_set("clientId", clientId)
memory_db_set("clientSecret", clientSecret)

add_header("Content-Type", "application/json")
add_header("X-ClientId", clientId)

return 200, "OK"

```

### Minimal routes.lua

`scripts/routes.lua`:
```lua
Routes = {
  { pattern = "*", response = { 200, "Hello, World!" } }
}
```


### Example before.lua

`scripts/before.lua`:
```lua
print("Before script")
add_header("X-Request-Id", "12345")
```

### Example after.Lua

`scripts/after.lua`:
```lua
add_header("Content-Type", "application/json")
```

## Avaliable Lua functions

- `add_header(name, value)`: Add a header to the response.
- `memory_db_set(key, value)`: Store a value in the in-memory database.
- `memory_db_get(key)`: Retrieve a value from the in-memory database.
- `memory_db_del(key)`: Delete a value from the in-memory database.
- `memory_db_flush()`: Clear all values from the in-memory database.
- `memory_db_keys()`: Retrieve a list of all keys in the in-memory database.
- `memory_db_values()`: Retrieve a list of all values in the in-memory database.
- `memory_db_pairs()`: Retrieve a list of all key-value pairs in the in-memory database.
- `memory_db_size()`: Retrieve the number of key-value pairs in the in-memory database.
- `memory_db_has(key)`: Check if a key exists in the in-memory database.
- `print(...)`: Print a message to the console.
- `sleep(seconds)`: Sleep for a specified number of seconds.


