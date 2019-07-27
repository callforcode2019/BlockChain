package main

import (
	"database/sql"
	"fmt"
	"github.com/BlockChain/blockchain-service/blockchain"
	"github.com/BlockChain/blockchain-service/web"
	"github.com/BlockChain/blockchain-service/web/controllers"
	"log"
	"os"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
	db,err := sql.Open("mysql","root:@/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	_,err = db.Exec("create table if not exists `user` (`username` varchar(20) not null primary key, `password` varchar(100) not null , `email` varchar(40) not null);")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("111111",os.Getenv("GOPATH"))
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
      	OrdererID: "orderer.hf.chainhero.io",
      	
		// Channel parameters
		ChannelID:     "chainhero",
		ChannelConfig: "/home/holdonbush/go/src/github.com/BlockChain/blockchain-service/fixtures/artifacts/chainhero.channel.tx",		// Chaincode parameters
		ChainCodeID:     "blockchain-service",
		ChaincodeGoPath: "/home/holdonbush/go",
		ChaincodePath:   "github.com/BlockChain/blockchain-service/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()	

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	// // Query the chaincode
	// response, err := fSetup.QueryHello()
	// if err != nil {
	// 	fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Response from the query hello: %s\n", response)
	// }

	// // Invoke the chaincode
	// txId, err := fSetup.InvokeHello("BlockChain")
	// if err != nil {
	// 	fmt.Printf("Unable to invoke hello on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Successfully invoke hello, transaction ID: %s\n", txId)
	// }

	// // Query again the chaincode
	// response, err = fSetup.QueryHello()
	// if err != nil {
	// 	fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Response from the query hello: %s\n", response)
	// }
	//Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
}