package main

import (
    "flag"
    "fmt"
    "bufio"
    "os"
    "log"
    "regexp"
)

type fileInfo struct {
    filepath string
    lines []string
    content string
    keywords []keywords
    filetype string
    
}
type keywords struct {
    keyword string
    line int
    matches []string
}
var verbose *bool
func main() {
    //fmt.Println(searchQuery("hack"))
    file := flag.String("file", "", "File you want to be searched")
    verbose = flag.Bool("verbose", false, "Verbose output")
    flag.Parse()
    if(*file == ""){
        panic("You need to provide a filepath. Use -h for help.")
    }
    var fileObj fileInfo
    fileObj.filepath = *file
    handleFile(&fileObj)
    fileObj.filetype = getFiletype(fileObj)[1]
    fmt.Printf("%#v\n", fileObj)
}

func doLog(text string){
    if(*verbose){
        //Log here
    }
}

func handleFile(fileObj *fileInfo){
    file, err := os.Open(fileObj.filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fileObj.lines = append(fileObj.lines, scanner.Text())
        fileObj.content = fileObj.content+"\n"+scanner.Text()
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
func getKeywords(fileObj fileInfo){
    var regex []string
    switch fileObj.filetype {
        case "py":
            regex = []string{
                "(\\w*?)\\W?=",
                "(([\"'])[^]*?\\2)",
                "(#.*?\n|\r)",
                "(\"\"\".*?\"\"\")",
            }

        case "php":
            regex = []string{
                "(\\$.*?)\\W?=",
                "(([\"'])[^]*?\\2)",
                "(\\/\\/.*?\\n|\\r)",
                "(\\/\\*.*?\\*\\/)",
            }
        case "go":
            regex = []string{
                "(([\"'])[^]*?\\2)",    
            }
    }
    for x := 0; x < len(regex); x++ {
        fmt.Println(regex[x])
    }
}
func getFiletype(fileObj fileInfo)([]string){
    regex := ".*?\\.([A-z]{2,3})$"
    r, _ := regexp.Compile(regex)
    match := r.FindStringSubmatch(fileObj.filepath)
    return match
}