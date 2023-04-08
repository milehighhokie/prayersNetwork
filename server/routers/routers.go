package routers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	// mySQL
	_ "github.com/go-sql-driver/mysql"
)

// Prayer = json struct to pass back via APIs
type Prayer struct {
	RowID           string `json:"idprayers"`
	PrayerText      string `json:"prayerText"`
	PrayerEntryDate string `json:"prayerEntryDate"`
	PrayerReadCount int32  `json:"prayerReadCount"`
}

// IncomingRecord
type IncomingRecord struct {
	IncomingString string `json:"prayerText"`
}

// IncomingUpdate
type IncomingUpdate struct {
	IncomingInt string `json:"prayerID"`
}

// RegisterRouters = serve up html and APIs
func RegisterRouters() *gin.Engine {
	var db *sql.DB
	userID := os.Getenv("PRAYERS_USER_ID")
	password := os.Getenv("PRAYERS_PASSWORD")
	connectString := userID + ":" + password + "@tcp(127.0.0.1:3306)/milehigh_prayers"
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()

	//
	// API section
	//

	// GET - Pull record info for all prayers
	router.GET("/api/record", func(c *gin.Context) {
		recordList, err := PrayerList(db)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, recordList)
	})

	//
	// POST - Create a record
	router.POST("/api/record", func(c *gin.Context) {
		var prayerText IncomingRecord

		var x []byte
		x, _ = ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(x, &prayerText)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("x is %s \n", x)
		fmt.Printf("prayerText is %s \n", prayerText.IncomingString)

		err = CreatePrayer(prayerText.IncomingString, db)
		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, x)
	})

	// PUT - update records
	router.PUT("/api/record", func(c *gin.Context) {
		var x []byte
		x, _ = ioutil.ReadAll(c.Request.Body)
		var prayerID IncomingUpdate
		err := json.Unmarshal(x, &prayerID)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("x is %s \n", x)
		fmt.Println(prayerID)

		recint, _ := strconv.Atoi(prayerID.IncomingInt)
		_, err = UpdatePrayer(recint, db)
		if err != nil {
			fmt.Println(err.Error())
		}

		c.JSON(http.StatusOK, x)

	})

	//return to main program
	return router
}
