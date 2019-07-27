package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

type Blocks struct {
	BlocksData []Cnt
	Username   string
}



type Cnt struct {
	Key		  string
	Value     string
	Sender    string
	Receive   string
	DataCnt	  string
	TimeStamp string
	TransactionId string
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	helloValue, err := app.Fabric.QueryHello()
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	data := &struct {
		Hello string
	}{
		Hello: helloValue,
	}
	fmt.Println(app.Fabric.OrdererID)
	db, err := leveldb.OpenFile("db",nil)
	defer db.Close()
	if err != nil {
		fmt.Println("open db err",err)
	}
	err = db.Put([]byte(app.Fabric.OrdererID+helloValue),[]byte(helloValue),nil)
	if err != nil {
		fmt.Println("put data err",err)
	}
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key),string(value))
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println("iter err",err)
	}

	renderTemplate(w, r, "home.html", data)
}

func (app *Application) HomeHandler_1(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		http.Redirect(w,r,"/login",http.StatusFound)
		return
	} else {
		fmt.Println(cookie.Value)
	}

	db, err := leveldb.OpenFile("db",nil)
	defer db.Close()
	if err != nil {
		fmt.Println("open db err",err)
	}
	iter := db.NewIterator(nil,nil)
	var blocks Blocks
	blocks.Username = cookie.Value
	for iter.Next() {
		var cnt Cnt
		key := iter.Key()
		value := iter.Value()

		var storedata StroeData
		err := json.Unmarshal(value, &storedata)
		if err != nil {
			log.Fatal("error",err)
		}
		fmt.Println(len(string(key)))
		cnt.Key = string(key)
		cnt.Value = string(value)
		cnt.Sender = storedata.Sender
		cnt.Receive = storedata.Receiver
		cnt.DataCnt = storedata.Data
		cnt.TransactionId = storedata.Transaction
		cnt.TimeStamp = time.Unix(storedata.TimeStamp,0).Format("02/01/2006 15:04:05 PM")
		blocks.BlocksData = append(blocks.BlocksData,cnt)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println("iter err",err)
	}
	renderTemplate(w,r,"home.html",blocks)
}