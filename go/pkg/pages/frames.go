package pages

import (
	"embed"
	_ "embed"
)

//go:embed frames/*.txt
var frameDir embed.FS

var base = `.""---------"".
|             /""\
|            | _  |
|             / | |
|             |/  |
|             /  /
|            |  /
|            "t"
'"---------"'`
