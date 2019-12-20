package database

import "database/sql"

//根据用户名返回密码
func GetUserPassword(username string) string{
	log.Debugf("database: query password by %s", username)
	rows, err := db.Query("select password from player where player.name = ?", username)
	if err != nil {
		log.Errorln(err)
		return ""
	}
	var password sql.NullString
	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			log.Errorln(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Errorln(err)
	}
	if password.Valid {
		return password.String
	}
	return ""
}