package filelines

import (
	"bufio"
	"embed"
	"fmt"
	"strings"
)

func EmbdCallFileLinesFunc(f func(string, int, interface{}) (interface{}, error), files *embed.FS, path, file string, src interface{}) error {
	byte, err := files.ReadFile(path + "/" + file)
	if err != nil {
		return fmt.Errorf("filelines.EmbdCallFileLinesFunc:embed.FS.ReadFile failure path:%q file:%q.\n\t%v", path, file, err)
	}
	var text string
	index := 0
	scanner := bufio.NewScanner(strings.NewReader(string(byte)))
	for scanner.Scan() {
		text = scanner.Text()
		if len(text) > 0 { // 空白行（末尾など）なら参照しない
			if text[0] != '#' { // 行頭が#ならコメント行なので参照しない
				if _, err = f(text, index, src); err != nil {
					return fmt.Errorf("filelines.EmbdCallFileLinesFunc:f failure text:%q.\n\t%v", text, err)
				}
				index++
			}
		}
	}
	return nil
}
