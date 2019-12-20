package database

import "database/sql"

const (
	UPDATE_SUCCESS = 300
	UPDATE_SERVER_ERROR = 301
	UPDATE_PASSWORD_ERROR = 302
)

func Win(username, password string, money int32) int32 {
	log.Debugf("user: %s win %s money\n", username, money)
	password1 := getUserPassword(username)
	if password1 != "" && password1 == password {
		res := pstx("update player as p set p.win = p.win + 1, p.money = p.money + ? where p.name = ?", money, username)
		num, err := res.RowsAffected()
		if err != nil {
			log.Errorln(err)
			return UPDATE_SERVER_ERROR
		}
		if num != 0 {
			return UPDATE_SUCCESS
		}
	}
	return UPDATE_PASSWORD_ERROR
}

func Lose(username, password string, money int32) int32 {
	log.Debugf("user: %s lose %s money\n", username, money)
	password1 := getUserPassword(username)
	if password1 != "" && password1 == password {
		res := pstx("update player as p set p.lose = p.lose + 1, p.money = p.money - ? where p.name = ?", money, username)
		num, err := res.RowsAffected()
		if err != nil {
			log.Errorln(err)
			return UPDATE_SERVER_ERROR
		}
		if num != 0 {
			return UPDATE_SUCCESS
		}
	}
	return UPDATE_PASSWORD_ERROR
}

func pstx(sql string, param ...interface{}) sql.Result{
	tx, err := db.Begin()
	if err != nil {
		log.Errorln(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Errorln(err)
	}
	res, err := stmt.Exec(param)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(res)
	err = tx.Commit()
	if err != nil {
		log.Errorln(err)
	}
	_ = stmt.Close()
	return res
}