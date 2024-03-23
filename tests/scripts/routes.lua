
Routers = {
  { pattern = "/services/oauth2/token", script = "oauth.lua",            methods = { "POST", "PUT", "PATCH" },  priority = 1 },
  { pattern = "*",                      response = { 404, "Not found" }, methods = { "GET" },                   priority = 100 },
  { pattern = "*",                      response = { 201, "" },          methods = { "POST", "PUT", "DELETE" }, priority = 101 },
  { pattern = "/",                      response = { 200, "" },          methods = { "ALL" },                   priority = 2 },
  { pattern = "/article/*",             response = { 201, "okey" },      methods = { "ALL" },                   priority = 3 }
}

return Routers
