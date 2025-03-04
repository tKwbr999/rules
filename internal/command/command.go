package command

import (
	"flag"
)

// コマンドライン引数を解析する
func ParseArgs() (string, []string, error) {
	editor := ""
	var files []string

	flag.Parse()

	if flag.NArg() > 0 {
		editor = flag.Arg(0)
	}

	for i := 1; i < flag.NArg(); i++ {
		files = append(files, flag.Arg(i))
	}

	return editor, files, nil
}
