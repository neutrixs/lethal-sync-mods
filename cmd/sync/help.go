package main

const help string =
"Usage: ./sync [options]\n"+
"\n"+
"Options:\n"+
"   -t <type>: Specifies the type of data to synchronize\n"+
"   <type>:\n"+
"       `mods`: Synchronizes mods (default)\n"+
"       `save`: Synchronizes savegame\n"+
"\n"+
"   --to-server: Instead of syncing to the client, it syncs the client to the server\n"+
"\n"+
"   -d <directory>: Specifies the directory of the client to sync to/from.\n"+
"                   Defaults to the working directory for type mods, and\n"+
"                   ..LocalLow/ZeekerssRBLX/Lethal Company for type save\n"+
"\n"+
"   --base-url <url>: Specifies base URL of the API\n"+
"                     defaults to https://lc.neutrixs.my.id\n"+
"\n"+
"   -h: Show this help text"