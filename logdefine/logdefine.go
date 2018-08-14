package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jslyzt/logdefine"
)

func main() {
	fileName := flag.String("file", "template.xml", "please input the xml file to analysis")
	outDir := flag.String("odir", "./out", "please input the dir out file to store")
	fileDir := flag.String("idir", "", "please input the dir files to analysis")
	outModel := flag.String("model", "go;cpp;java", "please input the file type to export")
	outCSet := flag.String("charset", "utf-8", "please input the output charset")
	incs := flag.String("inc", "", "please input the include files")
	flag.Parse()

	if (fileName == nil || len(*fileName) <= 0) && (fileDir == nil || len(*fileDir) <= 0) {
		fmt.Println("file or idir should set one")
		os.Exit(0)
	}
	//fmt.Println(*fileName, *outDir, *fileDir, *outModel)
	exportModel := make([]int8, 0)
	if outModel == nil || *outModel == "go;cpp;java" {
		exportModel = []int8{
			logdefine.ETgo,
			logdefine.ETcpp,
			logdefine.ETjava,
		}
	} else {
		for _, model := range strings.Split(*outModel, ";") {
			switch model {
			case "go":
				exportModel = append(exportModel, logdefine.ETgo)
			case "cpp":
				exportModel = append(exportModel, logdefine.ETcpp)
			case "java":
				exportModel = append(exportModel, logdefine.ETjava)
			}
		}
	}
	appends := make(map[string]interface{})
	if incs != nil && len(*incs) > 0 {
		appends["includes"] = strings.Split(*incs, ";")
	}
	if outCSet != nil {
		appends["charset"] = *outCSet
	}

	if fileName != nil && len(*fileName) > 0 {
		logfile := logdefine.AnalysisFile(*fileName)
		if logfile != nil {
			logfile.Export(exportModel, *outDir, appends)
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
				logfile := logdefine.AnalysisFile(file.Name())
				if logfile != nil {
					logfile.Export(exportModel, *outDir, appends)
				}
			}
		}
	}
}
