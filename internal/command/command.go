package command

import (
	"flag"
)

// コマンドライン引数を解析する
func ParseArgs() (string, string, error) {
	editor := ""
	env := ""

	flag.Parse()

	if flag.NArg() > 0 {
		editor = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		env = flag.Arg(1)
	}

	return editor, env, nil
}
