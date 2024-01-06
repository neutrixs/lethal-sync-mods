package main

const help = 
"Usage: ./gen_checksums [options]\n"+
"\n"+
"Options:\n"+
"   -t <type>: Specifies the type of data to generate checksums of\n"+
"   <type>:\n"+
"       `mods`: Generates checksums of mods (default)\n"+
"       `save`: Generates checksums of savegames\n"+
"\n"+
"   -d <dir>: Specifies the directory of the data.\n"+
"             Defaults to {Working Directory} for `-t mods`,\n"+
"             and {Working Directory}/saves for `-t save`"+
"\n"+
"   -h: Show this help text"