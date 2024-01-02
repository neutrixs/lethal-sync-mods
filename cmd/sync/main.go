package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neutrixs/lethal-sync-mods/constants"
	"github.com/neutrixs/lethal-sync-mods/internal/api"
)

const baseURL = "https://lc.neutrixs.my.id"

func main() {
	wd, err := os.Getwd()
	if err != nil { log.Fatalln(err) }

	err = api.SyncToClient(baseURL, wd, constants.ModsWhitelist, constants.ModsIgnore)
	if err != nil { log.Fatalln(err) }

	fmt.Println("Done! Press the Enter Key to exit anytime")
    fmt.Scanln()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}