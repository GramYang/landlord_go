package database

import (
	"database/sql"
)

func getUserPassword(username string) string{
	rows, err := db.Query("select password from player where player.name = ?", username)
	defer rows.Close()
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

func GetUserInfo(username string) (map[string]string, error) {
	rows, err := db.Query("select * from player where player.name = ?", username)
	defer rows.Close()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	var id,name,password,avatar,win,lose,money sql.NullString
	for rows.Next() {
		err = rows.Scan(&id,&name,&password,&avatar,&win,&lose,&money)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	result := make(map[string]string)
	result["id"] = id.String
	result["name"] = name.String
	result["password"] = password.String
	result["avatar"] = avatar.String
	result["win"] = win.String
	result["lose"] = lose.String
	result["money"] = money.String
	return result, nil
}