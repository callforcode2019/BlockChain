package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
	"time"
)

type BlockData struct {
	Sender      string
	Receiver    string
	Data		string
	TimeStamp   int64
}
type StroeData struct {
	Sender      string
	Receiver    string
	Data		string
	TimeStamp   int64
	Transaction string
}

func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		http.Redirect(w,r,"/login",http.StatusFound)
		return
	}
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
		Username      string
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
		Username:      cookie.Value,
	}

	if r.FormValue("submitted") == "true" {
		helloValue := r.FormValue("hello")
		receiver := r.FormValue("receive")
		data1 := r.FormValue("data")
		if receiver==""||data1=="" {
			fmt.Println("get content error")
		}
		var blockdata BlockData
		var storetodb StroeData
		blockdata.Sender = helloValue
		blockdata.Receiver = receiver
		blockdata.Data = data1
		blockdata.TimeStamp = time.Now().Unix()



		d, err := json.Marshal(&blockdata)
		if err != nil {
			fmt.Println("change to json err")
		} else {
			fmt.Println(string(d))
		}
		storedata := string(d)
		txid, err := app.Fabric.InvokeHello(storedata)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		storetodb.Sender = helloValue
		storetodb.Receiver = receiver
		storetodb.Data = data1
		storetodb.TimeStamp = blockdata.TimeStamp
		storetodb.Transaction = txid

		e, err := json.Marshal(&storetodb)
		if err != nil {
			fmt.Println("change to json err2")
		} else {
			fmt.Println(string(e))
		}
		storetodbdata := string(e)
		db, err := leveldb.OpenFile("db",nil)
		defer db.Close()
		if err != nil {
			fmt.Println("open db err",err)
		}
		err = db.Put([]byte(strconv.FormatFloat(float64(1)/float64(storetodb.TimeStamp),'E',-1,64)),[]byte(storetodbdata),nil)
		if err != nil {
			fmt.Println("Put data err",err)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "request.html", data)
}