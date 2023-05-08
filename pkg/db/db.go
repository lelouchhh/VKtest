package db

import (
	"VKtest/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Data struct {
	Service  string `db:"service"`
	Password string `db:"password"`
	Login    string `db:"service_login"`
	Time     string `db:"time"`
}

func GetDb(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Dbname, c.SslMode))
	if err != nil {
		return nil, err
	}
	return db, err
}

func RegisterUser(db *sql.DB, id string) (string, error) {
	//check if exist
	var count int
	rows, err := db.Query("SELECT count(*) FROM main.user WHERE LOGIN=$1", id)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return "Внутренняя ошибка", err
		}
	}
	if count > 0 {
		return "Вы уже зарегистрированы", nil
	}
	_, err = db.Exec("insert into main.user (LOGIN) VALUES ($1);", id)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	return "Вы успешно зарегистрированы", nil
}

func DeleteUserData(db *sql.DB, id, service, login string) (string, error) {
	var count int
	rows, err := db.Query(
		"select count(id) from main.data where service_login = $1 and service=$2 and login='$3;",
		login,
		service,
		id,
	)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return "Внутренняя ошибка", err
		}
	}
	if count == 0 {
		return "Запись не найдена", nil
	}
	_, err = db.Exec(
		"DELETE FROM main.data WHERE service_login=$1 and service=$2 and login=$3;",
		login,
		service,
		id,
	)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	return "Запись успешно удалена", nil
}

func GetUserData(db *sql.DB, id, service, login string) ([]Data, error) {
	var output []Data
	rows, err := db.Query("select service, service_login, password, end_time-now() as time from main.data where service=$1 and service_login = $2 and login=$3 and extract(epoch from (end_time - now())) > 0;",
		service,
		login,
		id)
	if err != nil {
		return []Data{}, err
	}
	for rows.Next() {
		var row Data
		err := rows.Scan(&row.Service, &row.Login, &row.Password, &row.Time)
		if err != nil {
			return nil, err
		}
		output = append(output, row)
	}

	return output, nil
}

func AddUserData(db *sql.DB, id, service, login, password string) (string, error) {
	var count int
	rows, err := db.Query("SELECT count(*) FROM main.data WHERE SERVICE_LOGIN=$1 and service=$2 and login=$3", login, service, id)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return "Внутренняя ошибка", err
		}
	}
	if count > 0 {
		_, err := db.Exec(fmt.Sprintf(
			"UPDATE main.data SET password = '%s',start_time='%s', end_time='%s' WHERE login = '%s' and service='%s' and service_login='%s';",
			password,
			time.Now().Format(time.DateTime),
			time.Now().Add(time.Hour*24*7).Format(time.DateTime),
			id,
			service,
			login))
		if err != nil {
			return "Внутренняя ошибка", err
		}
		return "Запись успешно обновлена", nil
	}
	_, err = db.Exec(
		"insert into main.data (login, service_login, service, password, start_time, end_time) values ($1,$2, $3, $4, $5, $6);",
		id,
		login,
		service,
		password,
		time.Now().Format(time.DateTime),
		time.Now().Add(time.Hour*24*7).Format(time.DateTime),
	)
	if err != nil {
		return "Внутренняя ошибка", err
	}
	return "Запись успешно добавлена", nil

}
