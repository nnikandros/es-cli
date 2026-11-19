/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "embed"

	"escobra/cmd"
)

//go:embed es_fields.yaml
var fieldsYAML []byte

func main() {

	cmd.Execute(fieldsYAML)
}
