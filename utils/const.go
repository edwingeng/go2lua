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
	{{- $n := sn .}}
local init_{{printf "%03d" $n}} = require("{{.}}")
if type(init_{{printf "%03d" $n}}) == "function" then
    init_{{printf "%03d" $n}}()
end
{{- end}}

return gopkg
`
)
