package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jslyzt/logDefine"
)

func main() {
	fileName := flag.String("file", "", "please input the xml file to analysis")
	outDir := flag.String("odir", "./out", "please input the dir out file to store")
	fileDir := flag.String("idir", "./", "please input the dir files to analysis")
	outModel := flag.String("model", "java", "please input the file type to export")
	flag.Parse()

	if (fileName == nil || len(*fileName) <= 0) && (fileDir == nil || len(*fileDir) <= 0) {
		fmt.Println("file or idir should set one")
		os.Exit(0)
	}
	exportModel := make([]int8, 0)
	if outModel == nil || *outModel == "go;cpp;java" {
		exportModel = logDefine.DefaultExport()
	} else {
		for _, model := range strings.Split(*outModel, ";") {
			switch model {
			case "go":
				exportModel = append(exportModel, logDefine.ET_GO)
			case "cpp":
				exportModel = append(exportModel, logDefine.ET_CPP)
			case "java":
				exportModel = append(exportModel, logDefine.ET_JAVA)
			}
		}
	}

	if fileName != nil && len(*fileName) > 0 {
		logfile := logDefine.AnalysisFile(*fileName)
		if logfile != nil {
			logfile.Export(exportModel, *outDir)
		}
	}
	if fileDir != nil && len(*fileDir) > 0 {
		dir, err := ioutil.ReadDir(*fileDir)
		if err != nil {
			fmt.Printf("read dir %s error %v\n", *fileDir, err)
			os.Exit(0)
		}
		for _, file := range dir {
			if file.IsDir() == false && strings.HasSuffix(strings.ToLower(file.Name()), ".xml") == true {
				logfile := logDefine.AnalysisFile(file.Name())
				if logfile != nil {
					logfile.Export(exportModel, *outDir)
				}
			}
		}
	}
}
