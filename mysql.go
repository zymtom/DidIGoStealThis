package main

import (
    //"os"
    //"fmt"
    //"log"
    "time"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)
/*
func main(){
    db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "", "didigostealthis")
    err := db.Connect()
    if err != nil {
        log.Fatal(err)
    }
    var fileObj fileInfo
    fileObj.filepath = "gomeme.go"
    fileObj.filetype = "go"
    fmt.Println(insertReport(db, fileObj))
    //rows, res, err := db.Query("select * from X where id > %d", 20)
}*/
func insertReport(db mysql.Conn, fileObj fileInfo)(error){
    stmt, err := db.Prepare("insert into reports (filepath, filetype, timestamp) values (?, ?, ?)")
    if err != nil {
        return err
    }
    _, err = stmt.Run(fileObj.filepath, fileObj.filetype, time.Now().Unix())
    return err
}
