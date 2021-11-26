package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	techStack *string
	sanitizedStack string
	wordlistFolder *string
	filePaths []string
	tmpFile *os.File
	stackItems []string
	wordlists []string
)

func main() {
	techStack = flag.String("techStack", "", "comma or space separated technologies (Apache,Microsoft Azure)")
	wordlistFolder = flag.String("wordlistFolder", "", "path to folder containing wordlist files named after technologies")
	flag.Parse()

	if *techStack == "" || *wordlistFolder == ""{
		os.Stderr.WriteString("techStack and wordlistFolder are required")
		return
	}

	if _, err := os.Stat(*wordlistFolder); err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	files, err := readFilesRecursive(*wordlistFolder)
	if err != nil {
		fmt.Println("Error reading only files from wordlistFolder: "+err.Error())
		return
	}
	for _,file := range files{
		filePaths = append(filePaths, file)
	}
	//tmpFile, err := os.CreateTemp("", "w0_")
	//defer tmpFile.Close()
	// sanitize techStack
	sanitizedStack = strings.Replace(*techStack," ",",",-1)
	sanitizedStack = strings.Replace(sanitizedStack,"-",",",-1)
	sanitizedStack = strings.Replace(sanitizedStack,"_",",",-1)
	sanitizedStack = strings.Replace(sanitizedStack,"(","",-1)
	sanitizedStack = strings.Replace(sanitizedStack,")","",-1)
	stackItems = RemoveEmpty(strings.Split(sanitizedStack, ","))

	// unique stackItems
	stackItems = removeDuplicateStr(stackItems)

	for _,path := range filePaths{
		for _,tech := range stackItems{
			if strings.Contains(strings.ToLower(strings.TrimSuffix(path,".txt")),strings.ToLower(tech)){
				if !stringInSlice(path,wordlists){
					wordlists = append(wordlists, path)
					file, _ := os.OpenFile(filepath.Join(*wordlistFolder,path),os.O_RDONLY,0644)
					content, err := io.ReadAll(file)
					if err != nil {
						fmt.Println("Error reading from file: "+err.Error())
						return
					}
					fmt.Println(string(content))
					//// seek to end of file for appending
					//tmpFile.Seek(0,2)
					//// append wordlist to tempfile
					//tmpFile.Write(content)
					file.Close()
				}
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func readFilesRecursive(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path,
		func(absPath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir(){
				files = append(files, strings.TrimPrefix(absPath,path+"/"))
			}
			d.Info()
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files, nil
}

func RemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" && str != " " {
			r = append(r, strings.TrimSpace(str))
		}
	}
	return r
}