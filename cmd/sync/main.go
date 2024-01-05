package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/neutrixs/lethal-sync-mods/constants"
	"github.com/neutrixs/lethal-sync-mods/internal/api"
)

const baseURL = "https://lc.neutrixs.my.id"
const savePath = "AppData/LocalLow/ZeekerssRBLX/Lethal Company"

func main() {
    wd, err := os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }
    homedir, err := os.UserHomeDir()
    if err != nil {
        log.Fatalln(err)
    }

    helpflag := flag.Bool("h", false, "show help")
    dirflag := flag.String("d", "", "set directory")
    typeflag := flag.String("t", "mods", "set type")
    baseurlflag := flag.String("base-url", baseURL, "set base URL")
    toserverflag := flag.Bool("to-server", false, "sync to server")
    userflag := flag.String("u", "", "specify user")

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
        dir := *dirflag
        if dir == "" {
            dir = wd
        }

        err = api.SyncToClient(
            *baseurlflag, dir, constants.ModsWhitelist, constants.ModsIgnore,
        )
        if err != nil { log.Fatalln(err) }
        fmt.Println("Done! Press the Enter Key to exit anytime")
        fmt.Scanln()
    }
    case "save":
        if *userflag == "" {
            fmt.Printf("please specify the user!\n%s\n", help)
            os.Exit(0)
        }

        dir := *dirflag
        if dir == "" {
            dir = path.Join(homedir, savePath)
        }

        parsed, err := url.Parse(*baseurlflag)
        if err != nil { log.Fatalln(err) }

        parsed.Path = path.Join(parsed.Path, "saves", *userflag)
        fmt.Println(parsed.String())

        err = api.SyncToClient(
            parsed.String(), dir, constants.SaveWhitelist, constants.SaveIgnore,
        )

        if err != nil { log.Fatalln(err) }
        fmt.Println("Done! Press the Enter Key to exit anytime")
        fmt.Scanln()
    default:
        fmt.Printf("invalid type option!\n%s\n", help)
    }
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}