package routers

import (
	"database/sql"
	"fmt"

	// mySQL
	_ "github.com/go-sql-driver/mysql"
)

// UpdatePrayer = Update prayers for idprayers
func UpdatePrayer(prayerUpdate int, db *sql.DB) (updated int64, err error) {
	sql := `
	      UPDATE prayers
			SET
			prayerReadCount = prayerReadCount + 1
			WHERE idprayers = ?
	     `

	stmt, err := db.Prepare(sql)
	if err != nil {
		err = fmt.Errorf("error in UpdatePrayer after Prepare: %s", err.Error())
		return 0, err
	}

	result, err := stmt.Exec(prayerUpdate)
	updated, _ = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("error in UpdatePrayer after Exec: %s", err.Error())
	}

	stmt.Close()
	return
}

// CreatePrayer = Create a new prayer
func CreatePrayer(prayer string, db *sql.DB) (err error) {

	stmt, err := db.Prepare("insert into prayers (prayerText, prayerEntryDate, prayerReadCount) values (?, CURDATE(), 0);")
	if err != nil {
		err = fmt.Errorf("error in CreatePrayer after Prepare: %s", err.Error())
	} else {
		result, err := stmt.Exec(prayer)
		inserted, _ := result.RowsAffected()
		if err != nil || inserted != 1 {
			err = fmt.Errorf("error in CreatePrayer after Exec: %s", err.Error())
		}
		stmt.Close()
	}
	return
}

// PrayerList = Get multiple prayers
func PrayerList(db *sql.DB) (prayerList []Prayer, err error) {
	var nextPrayer Prayer
	var idprayers, prayerText, prayerEntryDate string
	var prayerReadCount int32
	rows, err := db.Query("select idprayers, prayerText, prayerEntryDate, prayerReadCount from prayers order by prayerReadCount asc, prayerEntryDate desc limit 200")
	if err != nil {
		err = fmt.Errorf("error in PrayerList after Query: %s", err.Error())
	} else {
		for rows.Next() {
			err = rows.Scan(&idprayers, &prayerText, &prayerEntryDate, &prayerReadCount)
			if err != nil {
				err = fmt.Errorf("error in PrayerList after Scan: %s", err.Error())
			} else {
				nextPrayer.RowID = idprayers
				nextPrayer.PrayerText = prayerText
				nextPrayer.PrayerEntryDate = prayerEntryDate
				nextPrayer.PrayerReadCount = prayerReadCount
				prayerList = append(prayerList, nextPrayer)
			}
		}
		rows.Close()
	}

	return
}
