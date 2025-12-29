package filelines

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"
)

func CallFileLinesFunc(f func(string, string, int, interface{}) error, path, file string, src interface{}) (int, error) {
	index := 0
	if fp, err := os.Open(path + "/" + file); err == nil {
		defer fp.Close()
		var line string
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line = scanner.Text()
			if len(line) > 0 { // 空白行（末尾など）なら参照しない
				if line[0] != '#' { // 行頭が#ならコメント行なので参照しない
					if err = f(line, path, index, src); err != nil {
						return 0, fmt.Errorf("filelines.CallFileLinesFunc:f failure line:%q.\n\t%v", line, err)
					}
					index++
				}
			}
		}

	} else {
		return 0, fmt.Errorf("filelines.CallFileLinesFunc:os.Open failure path:%q file:%q.\n\t%v", path, file, err)
	}
	return index, nil
}

func EmbdCallFileLinesFunc(f func(string, string, int, interface{}) error, files *embed.FS, path, file string, src interface{}) (int, error) {
	byte, err := files.ReadFile(path + "/" + file)
	if err != nil {
		return 0, fmt.Errorf("filelines.EmbdCallFileLinesFunc:embed.FS.ReadFile failure path:%q file:%q.\n\t%v", path, file, err)
	}
	var line string
	index := 0
	scanner := bufio.NewScanner(strings.NewReader(string(byte)))
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > 0 { // 空白行（末尾など）なら参照しない
			if line[0] != '#' { // 行頭が#ならコメント行なので参照しない
				if err = f(line, path, index, src); err != nil {
					return index, fmt.Errorf("filelines.EmbdCallFileLinesFunc:f failure line:%q.\n\t%v", line, err)
				}
				index++
			}
		}
	}
	return index, nil
}
