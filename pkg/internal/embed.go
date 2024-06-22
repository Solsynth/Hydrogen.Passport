package pkg

import "embed"

//go:embed all:views/*
var FS embed.FS
