package utils

const (
	GopkgFile = `
local gopkg = _G["%s"]
do
    local g = _G
    local newEnv = setmetatable({}, {
        __index = function (t, k)
            local v = gopkg[k]
            if v == nil then return g[k] end
            return v
        end,
        __newindex = gopkg,
    })
    _ENV = newEnv
end
`
)

const (
	TemplateGopkg = `{{"" -}}
-- package: {{.PkgName}}

local gopkg = {}
_G["{{.PkgPath}}"] = gopkg
do
    local g = _G
    local newEnv = setmetatable({}, {
        __index = function (t, k)
            local v = gopkg[k]
            if v == nil then return g[k] end
            return v
        end,
        __newindex = gopkg,
    })
    _ENV = newEnv
end

local checkModName = function (name)
	if package.loaded[name] ~= nil then
		error("mod name collision detected. name: " .. name)
	end
end
{{""}}
{{- range .Files}}
checkModName("{{.}}")
{{- end}}
{{""}}
local hashes = {
{{- range .Files}}
	{{- $h := hash .}}
	init_{{$h}} = true,
{{- end}}
}
local hashCount = 0
for _ in pairs(hashes) do
	hashCount = hashCount + 1
end
if hashCount ~= {{len .Files}} then
	error("hash collision detected")
end
{{""}}
{{- range .Files}}
	{{- $h := hash .}}
local init_{{$h}} = require("{{.}}")
{{- end}}

{{- if .Initializers}}

-- Initializers
{{.Initializers}}
{{- end}}

{{- range .Files}}
	{{- $h := hash .}}
init_{{$h}}() -- {{.}}
{{- end}}

return gopkg
`
)
