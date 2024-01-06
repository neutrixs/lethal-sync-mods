package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/neutrixs/lethal-sync-mods/constants"
)

func main() {
    wd, err := os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }

    helpflag := flag.Bool("h", false, "show help")
    typeflag := flag.String("t", "mods", "specifies type")
    dirflag := flag.String("d", "", "specifies directory")

    flag.Parse()

    if *helpflag {
        fmt.Println(help)
        os.Exit(0)
    }

    switch strings.ToLower(*typeflag) {
    case "mods": {
        dir := *dirflag
        if dir == "" {
            dir = wd
        }

        err = Generate(dir, constants.ModsWhitelist, constants.ModsIgnore)
        if err != nil { log.Fatalln(err) }
    }
    case "save": {
        dir := *dirflag
        if dir == "" {
            dir = path.Join(wd, "saves")
        }

        d, err := os.Open(dir)
        if err != nil { log.Fatalln(err) }

        stat, err := d.Stat()
        if err != nil { log.Fatalln(err) }

        if !stat.IsDir() {
            log.Fatalf("%s is NOT a directory!\n", dir)
        }

        children, err := d.Readdir(-1)
        if err != nil { log.Fatalln(err) }

        var saves []string

        for _, child := range children {
            if !child.IsDir() { continue }
            fullpath := path.Join(dir, child.Name())

            saves = append(saves, fullpath)
        }

        for _, save := range saves {
            err := Generate(save, constants.SaveWhitelist, constants.SaveIgnore)
            if err != nil { log.Fatalln(err) }
        }
    }
    default:
        fmt.Printf("invalid type option!\n%s\n", help)
        os.Exit(1)
    }
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}