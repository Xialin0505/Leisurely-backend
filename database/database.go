package database

import (

	//"os"

	"encoding/json"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"os"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ENVSETTING struct {
	BackendServer_Name string `json:"BackendServer_Name"`
	BackendRecom_Name  string `json:"BackendRecom_Name"`
	Database_Name      string `json:"Database_Name"`
}

var (
	DB          *gorm.DB
	ENVSettings ENVSETTING
)

func fetchEnvSetting() ENVSETTING {
	jsonFile, err := os.Open("./setting.json")
	if err != nil {
		fmt.Println(err)
	}

	var settings ENVSETTING

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &settings)

	if err != nil {
		println("error when fetching pop-up event")

	}
	return settings
}

func init() {

	settings := fetchEnvSetting()
	ENVSettings = settings

	//hostname, err := os.Hostname()
	//dbn := fmt.Sprintf("host=192.168.2.105 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	dbn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=postgres port=5432 sslmode=disable", ENVSettings.Database_Name)
	//hostname, os.Getenv("POSTGRES_PASSWORD"))

	db, err := gorm.Open(postgres.Open(dbn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	DB = db
}
