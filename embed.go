package mailtplsrv

import "embed"

//go:embed tmpl/*
var Templates embed.FS

//go:embed static/*
var Static embed.FS
