package constants

const SavesLocation = "AppData/LocalLow/ZeekerssRBLX/Lethal Company"

var ModsWhitelist = []string{
    "winhttp.dll",
    "doorstop_config.ini",
    "BepInEx/*",
}

var ModsIgnore = []string{
    "BepInEx/LogOutput.log",
    "BepInEx/cache/*",
}

var SaveWhitelist = []string{
    "*",
}

var SaveIgnore = []string{
    "checksums.txt",
}