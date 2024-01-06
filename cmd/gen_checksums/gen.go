package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/neutrixs/lethal-sync-mods/internal/api"
)

func Generate(dir string, whitelist []string, ignore []string) error {
	cs, err := api.GetChecksums(dir, whitelist, ignore)
	if err != nil { return err }

	file, err := os.OpenFile(path.Join(dir, "checksums.txt"), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
	if err != nil { return err }
	defer file.Close()

	csbyte, err := json.Marshal(cs)
	if err != nil { return err }

	_, err = file.Write(csbyte)
	if err != nil { return err }

	return nil
}