package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/neutrixs/lethal-sync-mods/constants"
	"github.com/neutrixs/lethal-sync-mods/internal/api"
)

const baseURL = "https://lc.neutrixs.my.id"

func main() {
    wd, err := os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }

    helpflag := flag.Bool("h", false, "show help")
    dirflag := flag.String("d", wd, "set directory")
    typeflag := flag.String("t", "mods", "set type")
    baseurlflag := flag.String("base-url", baseURL, "set base URL")
    toserverflag := flag.Bool("to-server", false, "sync to server")

    flag.Parse()

    if *helpflag {
        fmt.Println(help)
        os.Exit(0)
    }

    if *toserverflag {
        fmt.Println("--to-server not yet implemented!")
        os.Exit(0)
    }

    switch *typeflag {
    case "mods": {
        err = api.SyncModsToClient(
            *baseurlflag, *dirflag, constants.ModsWhitelist, constants.ModsIgnore,
        )
        if err != nil { log.Fatalln(err) }
        fmt.Println("Done! Press the Enter Key to exit anytime")
        fmt.Scanln()
    }
    case "save":
        fmt.Println("-t save not yet implemented!")
    default:
        fmt.Printf("invalid type option!\n%s\n", help)
    }
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}