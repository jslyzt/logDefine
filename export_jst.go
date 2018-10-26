package logdefine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func jsfmortLogfile(file *XMLLogFile) []byte {
	return nil
}

// 导出 javascript
func (file *XMLLogFile) exportJS(outdir string) bool {
	data := jsfmortLogfile(file)
	if data != nil && len(data) > 0 {
		fileName := fmt.Sprintf("%s/%s.js", outdir, file.Name)
		fmt.Printf("save file: %s\n", fileName)

		os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
		err := ioutil.WriteFile(fileName, data, os.ModePerm)
		if err != nil {
			fmt.Printf("save file %s failed, error %v", fileName, err)
			return false
		}
	}
	return true
}
