package main

import (
    "flag"
    "fmt"
    "bufio"
    "os"
    "log"
    "regexp"
    "strings"
    "strconv"
    //"errors"
    //"time"
    //"sync"
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
var verbose bool
var logToFile bool
var threads int
var retry int
var sleep int
func main() {
    //fmt.Println(searchQuery("hack"))
    paramMap := map[string][]string{
        "file":[]string{"string", "", "File you want to be searched"},
        "file-extension":[]string{"string", "", "Extension for the file, if it has none. e.g 'go'"},
        "verbose":[]string{"bool", "false", "Verbose output"},
        "log":[]string{"bool", "false", "Turns on log, which by default logs into the standard file 'logs'"},
        "threads":[]string{"int", "2", "Number of threads to be used while downloading the data."},
        "retry":[]string{"int", "3", "Number of times you want the program to attempt to retry downloading data from the website"},
        "sleep":[]string{"int", "3", "Time to sleep between retrying to download data."},
    }
    //paramMap["file"] = []string{"string", "", "File you want to be searched"}
    values := handleParams(paramMap)
    if(values["file"] == ""){
        log.Fatal("You need to provide a filepath. Use -h for help.")
    }
    verbose = values["verbose"].(bool)
    logToFile = values["log"].(bool)
    threads = values["threads"].(int)
    retry = values["retry"].(int)
    sleep = values["sleep"].(int)
    var fileObj fileInfo
    fileObj.filepath = values["file"].(string)
    handleFile(&fileObj)
    if values["file-extension"] != "" {
        fileObj.filetype = values["file-extension"].(string)
    }else{
        getFiletype(&fileObj)
    }
    getKeywords(&fileObj)
    fmt.Printf("%#v\n", fileObj.keywords)
}

func doLog(text string){
    if(verbose){
        fmt.Println(text)
    }
    if(logToFile){
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
                "(\\/\\/(.*)(?:\n|$|\r))",
                
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
/*func search(fileObj *fileInfo){
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
}*/
func handleParams(params map[string][]string)(map[string]interface{}){
    args := map[string]interface{}{}
    args["config"] = flag.String("config", "", "Config file to read from, may be used as an alternative to writing cli over and over")
    for k, v := range params {
        if(len(v) != 3){
            log.Fatal("Not enough arguments for flag: "+k)
        }
        types := v[0]
        defaultvalue := v[1]
        text := v[2]
        
        if types == "string" || types == "str" {
            args[k] = flag.String(k, defaultvalue, text)  
        }else if types == "int" || types == "integer" {
            if i, err := strconv.Atoi(defaultvalue); err == nil {
                args[k] = flag.Int(k, i, text)
                //fmt.Println(reflect.TypeOf(meme))  
            }else{
                log.Fatal("Invalid default value for flag: "+k + "| "+err.Error())
            }
        }else if types == "bool" || types == "boolean" {
            if strings.ToLower(defaultvalue) == "true" {
                args[k] = flag.Bool(k, true, text) 
            }else if strings.ToLower(defaultvalue) == "false"{
                args[k] = flag.Bool(k, false, text)
            }else{
                log.Fatal("Invalid default value for flag: "+k)
            }
            
        }           
        
    }
    flag.Parse()
    config := make(map[string]string)
    /*fmt.Printf("%#v\n", args)
    if str, ok := args["config"].(*string); ok {
        fmt.Printf("%#v\n", *str)
        meme = *str
        fmt.Printf("%#v\n", meme)
    }*/

    if k, v := args["config"].(*string); v {
        if *k != ""{
            file, err := os.Open(*k)
            if err != nil {
                doLog("[-]Error retrieving content for config")
                log.Fatal(err)
            }
            defer file.Close()
            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
                ex := strings.Split(scanner.Text(), "=")
                for _, v := range ex[1:] {
                    config[ex[0]] = config[ex[0]]+v
                }
            }
            if err := scanner.Err(); err != nil {
                doLog("[-]Error retrieving content for config")
                log.Fatal(err)
            }
        }
    }
    flags := map[string]interface{}{}
    for k, v := range config {
        if strings.ToLower(v) == "true" {
            flags[k] = true
        }else if strings.ToLower(v) == "false"{
            flags[k] = false
        }else if i, err := strconv.Atoi(v); err == nil {
            flags[k] = i
        }else{
            flags[k] = v
        }
            
    }
    for k, v := range params {
        if _, vc := flags[k]; vc  {
            if str, ok := args[k].(*string); ok {
                if v[1] != *str {
                    flags[k] = *str
                }
            }else if str, ok := args[k].(*int); ok {
                incInterface := stringToValidType(v[1])
                if incInterface.(int) != *str {
                    flags[k] = *str
                }
            }else if str, ok := args[k].(*bool); ok {
                incInterface := stringToValidType(v[1])
                if incInterface.(bool) != *str {
                    flags[k] = *str
                }
            }
        }else{
            if str, ok := args[k].(*string); ok {
                flags[k] = *str
            }else if str, ok := args[k].(*int); ok {
                flags[k] = *str
            }else if str, ok := args[k].(*bool); ok {
                flags[k] = *str
            }
        }
    }
    return flags
    /*file := flag.String("file", "", "File you want to be searched")
    filetype := flag.String("file-extension", "", "Extension for the file, if it has none. e.g 'go'")
    verbose = flag.Bool("verbose", false, "Verbose output")
    logToFile = flag.Bool("log", false, "Turns on logging into the standard file 'logs'")
    threads = flag.Int("threads", 2, "Number of threads to be used while downloading the data.")
    retry = flag.Int("retry", 3, "Number of times you want the program to attempt to retry downloading data from the website")
    sleep = flag.Int("sleep", 3, "Time to sleep between retrying to download data")
    flag.Parse()*/
}
func stringToValidType(str string)(interface{}){
    if strings.ToLower(str) == "true" {
        return true
    }else if strings.ToLower(str) == "false"{
        return false
    }else if i, err := strconv.Atoi(str); err == nil {
        return i
    }else{
        return str
    }
}