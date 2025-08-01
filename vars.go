package vars

import (
	"LocalDex/types"
	"embed"
)

//go:embed client/dist/**
var ViteFS embed.FS

//go:embed assets/**
var AssetsFS embed.FS

var ReverseProxy types.BehindReverseProxy
