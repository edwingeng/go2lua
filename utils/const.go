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
{{""}}
{{- range .Files}}
	{{- $h := hash .}}
local init_{{$h}} = require("{{.}}")
if type(init_{{$h}}) == "function" then
    init_{{$h}}()
end
{{- end}}

return gopkg
`
)
