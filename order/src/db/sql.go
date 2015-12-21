package db

import (
	"log"

	"github.com/Dataman-Cloud/seckilling/order/src/model"
)

func InsertOrder(order model.Order) error {
	sql := `insert into order(eid, uid, status, seq, ext, create) values (:eid, :uid, :status, :seq, :ext, :create)`

	dbConn := DB()
	stmt, err := dbConn.PrepareNamed(sql)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println("close stmt has error: ", err)
		}
	}()

	result, err := stmt.Exec(order)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetEId() (int64, error) {
	db := DB()
	var eid int64
	err := db.Get(&eid, "select id from order")
	if err != nil {
		log.Println("Get eid has error: ", err)
	}

	return eid, err
}
