package main

import (
    "flag"
    "fmt"
    "bufio"
    "os"
    "log"
    "regexp"
    "strings"
)

type fileInfo struct {
    filepath string
    lines []string
    content string
    keywords []keyword
    filetype string
    
}
type keyword struct {
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
    getKeywords(&fileObj)
    fmt.Printf("%#v\n", fileObj.keywords)
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
func getKeywords(fileObj *fileInfo){
    var regex []string
    switch fileObj.filetype {
        case "py":
            regex = []string{
                "(\\w*?)\\W?=",
                "\"([\\s\\S]*?)\"",
                "'([\\s\\S]*?)'",
                "(#.*?\n|\r)",
                "(\"\"\".*?\"\"\")",
            }

        case "php":
            regex = []string{
                "(\\$.*?)\\W?=",
                "\"([\\s\\S]*?)\"",
                "'([\\s\\S]*?)'",
                "(\\/\\/.*?\\n|\\r)",
                "(\\/\\*.*?\\*\\/)",
            }
        case "go":
            regex = []string{
                "\"([\\s\\S]*?)\"",
                "'([\\s\\S]*?)'",
                "type ([A-z0-9]*?) struct",
                "var ([A-z0-9]*?) ",
                "func ([A-z0-9]*?)\\(",
                "([A-z])*? :=",
                "package ([A-z0-9]*?)\\n",
                "\\/\\*([\\S\\s]*)\\*\\/",
                "\\/\\/(.*)(?:\n|$)",
                
            }
    }
    for x := 0; x < 1; x++ {
        //r := regexp.MustCompile(regex[x])
        r, _ := regexp.Compile(regex[x])
        matches := r.FindAllStringSubmatch(fileObj.content, -1)
        for _, match := range matches {
            var keyword keyword
            for i := 0; i < len(fileObj.lines); i++{
                if strings.Contains(fileObj.lines[i], match[1]) {
                    keyword.line = i
                }
            }
            keyword.keyword = match[1]
            fileObj.keywords = append(fileObj.keywords, keyword)
        }
        fmt.Printf("%#v\n\n", matches[1])
    }
    //fileObj.keywords = append(fileObj.keywords, keywords)
}

func getFiletype(fileObj fileInfo)([]string){
    regex := ".*?\\.([A-z]{2,3})$"
    r, _ := regexp.Compile(regex)
    match := r.FindStringSubmatch(fileObj.filepath)
    return match
}
//meme