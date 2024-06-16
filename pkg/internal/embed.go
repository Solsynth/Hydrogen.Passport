package pkg

import "embed"

//go:embed views/*
var FS embed.FS
