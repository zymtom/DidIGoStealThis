package main

import (
    "flag"
    "fmt"
    "bufio"
    "os"
    "log"
    "regexp"
    "strings"
    "time"
)

type fileInfo struct {
    filepath string
    lines []string
    content string
    keywords []keyword
    filetype string
    urls []string
}
type keyword struct {
    keyword string
    line int
    matches []string
}
type urlContent struct {
    url string
    content string
}
var verbose *bool
var logToFile *bool
var threads *int
var retry *int
var sleep *int
func main() {
    //fmt.Println(searchQuery("hack"))
    file := flag.String("file", "", "File you want to be searched")
    filetype := flag.String("file-extension", "", "Extension for the file, if it has none. e.g 'go'")
    verbose = flag.Bool("verbose", false, "Verbose output")
    logToFile = flag.Bool("log", false, "Turns on logging into the standard file 'logs'")
    threads = flag.Int("threads", 2, "Number of threads to be used while downloading the data.")
    retry = flag.Int("retry", 3, "Number of times you want the program to attempt to retry downloading data from the website")
    sleep = flag.Int("sleep", 3, "Time to sleep between retrying to download data")
    flag.Parse()
    if(*file == ""){
        log.Fatal("You need to provide a filepath. Use -h for help.")
    }
    var fileObj fileInfo
    fileObj.filepath = *file
    handleFile(&fileObj)
    if *filetype != "" {
        fileObj.filetype = *filetype
    }else{
        getFiletype(&fileObj)
    }
    getKeywords(&fileObj)
    fmt.Printf("%#v\n", fileObj.keywords)
}

func doLog(text string){
    if(*verbose){
        fmt.Println(text)
    }
    if(*logToFile){
        file, err := os.Open("logs")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        
        f, err := os.Create("logs")
        w := bufio.NewWriter(f)
        _, err = w.WriteString(text+"\n")
        if err != nil {
            log.Fatal(err)
        }
        w.Flush()
    }
}

func handleFile(fileObj *fileInfo){
    file, err := os.Open(fileObj.filepath)
    if err != nil {
        doLog("[-]Error retrieving content for file")
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fileObj.lines = append(fileObj.lines, scanner.Text())
        fileObj.content = fileObj.content+"\n"+scanner.Text()
    }

    if err := scanner.Err(); err != nil {
        doLog("[-]Error retrieving content for file")
        log.Fatal(err)
    }
    doLog("[+]Retrieved content for file")
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
                "(\"([\\s\\S]*?)\")",
                "('([\\s\\S]*?)')",
                "(type ([A-z0-9]*?) struct)",
                "(var ([A-z0-9]*?) )",
                "(func ([A-z0-9]*?)\\()",
                "(([A-z])*? :=)",
                "(package ([A-z0-9]*?)\\n)",
                "(\\/\\*([\\S\\s]*)\\*\\/)",
                "(\\/\\/(.*)(?:\n|$))",
                
            }
    }
    for x := 0; x < 1; x++ {
        //r := regexp.MustCompile(regex[x])
        r, _ := regexp.Compile(regex[x])
        matches := r.FindAllStringSubmatch(fileObj.content, -1)
        //fmt.Printf("Found keywords: %#v\n", matches)
        doLog("[+]Found keywords")
        for _, match := range matches {
            var keyword keyword
            for i := 0; i < len(fileObj.lines); i++{
                if strings.Contains(fileObj.lines[i], match[1]) {
                    keyword.line = i
                    break
                }
            }
            keyword.keyword = match[2]
            fileObj.keywords = append(fileObj.keywords, keyword)
        }
        //fmt.Printf("%#v\n\n", matches[1])
    }
    //fileObj.keywords = append(fileObj.keywords, keywords)
}

func getFiletype(fileObj *fileInfo){
    regex := ".*?\\.([A-z]{2,3})$"
    r, _ := regexp.Compile(regex)
    match := r.FindStringSubmatch(fileObj.filepath)
    if match == nil {
        doLog("[-]The filetype could not be found.")
        log.Fatal("The filetype could not be found. Specify via argument.")
    }
    supportedLangs := []string{"go", "php", "py"}
    if !stringInSlice(match[1], supportedLangs) {
        doLog("[-]The filetype is not supported yet. Either you can make an issue on the github for this(http://github.com/zymtom/didigostealthis) or you can write your own regex's and make a pull request")
        log.Fatal("This filetype is not supported yet.")
    }
    fmt.Printf("%#v \n",match)
    doLog("[+]Found filetype: "+match[1])
    fileObj.filetype = match[1]
}
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
func search(fileObj *fileInfo){
    tasks := make(chan string)
    results := make(chan urlContent)
    var sync sync.WaitGroup
    for i := 0; i < threads; i++ {
        sync.Add(1)
        go func() {
            for website := range tasks {
                tries := 0
                for {
                    if tries == retry{
                        time.Sleep(sleep)
                        break
                    }
                    r := getUrlContent(website)
                    tries++
                    if r != nil {
                        var res urlContent
                        res.url = website
                        res.content = r
                        results <- res
                        break
                    }else{
                        //log
                    }
                }
            }
            sync.Done()
        }()
    }
    
    for i := 0; i < len(websites); i++ {
        tasks <- websites[i]
    }
    close(tasks)
    sync.Wait()
    uniqueUrls := make(map[string]string)
    for obj := range results {
    
    }
}
func getUrlContent(url string) string {
    res, err := http.Get(url)
    if err != nil { 
        return nil
    } else { 
        body, _ := ioutil.ReadAll(res.Body)
        res.Body.Close()
        conv := string(body[:])
        return conv
    }
}
//meme