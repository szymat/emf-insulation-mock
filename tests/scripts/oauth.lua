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

add_header("Content-Type", "application/json")
add_header("X-ClientId", clientId)
add_header("X-ClientSecret", clientSecret)

return 200, "OK"
