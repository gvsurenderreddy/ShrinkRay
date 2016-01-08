package main

import (
	"os"
	"fmt"
	"flag"
	"regexp"
	"io/ioutil"
	"path/filepath"
)

//todo: optional conf file have excluded directories and files. ignorable keywords from js, css, html or customizable framework commands that cannot be changed

type UniqueWord struct {
	count int
	name string
	newName string
	discoveredIn []string
}

func main() {

	flag.Parse()

	pathToScan := flag.Arg(0)

	//change to return array. only grab html js and css
	validFiles, err := GetFilesInDirectory(pathToScan)

	if err != nil {fatal("Problem reading file or folder", err)}

	ignoreFiles := map[string]bool{
		"W:\\work\\git\\CustomTaxApp\\www\\css\\kendo.mobile.all.min.css": true,
		"W:\\work\\git\\CustomTaxApp\\www\\js\\libs\\jquery.min.js": true,
		"W:\\work\\git\\CustomTaxApp\\www\\js\\libs\\kendo.all.min.js": true,
		"W:\\work\\git\\CustomTaxApp\\www\\js\\libs\\kendo.autocomplete.min.js": true,
	}

	uniqueWords := make(map[string]UniqueWord, 0)

	for _, fileName := range validFiles{

		_, isIgnore := ignoreFiles[fileName]

		if(!isIgnore){

			fmt.Println("Examining: " + fileName)

			buf, err := ioutil.ReadFile(fileName)

			if err != nil {fatal("File exists but cannot be read?", err)}

			//now regex each file and store all the unique words found.
			words := regexp.MustCompile(`[A-Za-z]\w+`).FindAllStringSubmatch(string(buf), -1)

			//fmt.Printf("%v \n\n\n %v \n\n\n\n\n", string(buf), words)

			for _, word := range words {

				firstWord:= word[0]

				_, isAlreadyFound := uniqueWords[firstWord]

				if(isAlreadyFound){

					uniqueWords[firstWord] = UniqueWord{
						count: (uniqueWords[firstWord].count+1),
						name: firstWord,
						newName: "todo",
						discoveredIn: []string{fileName},
					}
					//todo, append discovered in.. meah
				}else if(!IsReservedWord(firstWord)){
					uniqueWords[firstWord] = UniqueWord{
						count: 1,
						name: firstWord,
						newName: "todo",
						discoveredIn: []string{fileName},
					}
				}

			}
		}
	}

	for _, uniqueWord := range uniqueWords{
		fmt.Printf("%v: %s\n", uniqueWord.count, uniqueWord.name)
	}
}

func GetFilesInDirectory(dir string, )([]string, error){

	files := make([]string, 0)

	walkErr := filepath.Walk(dir, func (objName string, info os.FileInfo, inErr error) error {

		if(!info.IsDir() && regexp.MustCompile(`(?i)^.*\.(js|html|css)$`).MatchString(objName)){
			files = append(files, objName)}

		return inErr
	})

	return files, walkErr
}


func IsReservedWord(someWord string)(bool){
	
	// the array of strings can come from the config file later on

	joinedReservations := []string{`and`, `iPhone`}

	for _, forbidden := range joinedReservations{
		if(someWord == forbidden){
			return true;
		}
	}

	return false
}

func fatal(myDescription string, err error){
	fmt.Println("FATAL: ", myDescription)
	fmt.Println(err)
	os.Exit(1)
}