package tools

import (
	"VKtest/pkg/db"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func Hash(s string) string {
	salted := []byte(s)
	salt := []byte("L0qu-gd0yr3D9*")
	salted = append(salted, salt...)
	h := sha1.New()
	h.Write(salted)
	hash := hex.EncodeToString(h.Sum(nil))

	return hash

}

func Prettify(data []db.Data) string {
	var output string
	if len(data) == 0 {
		return "Данных по данному запросу не найдено"
	}
	for _, val := range data {
		output += fmt.Sprintf("------------------------------\nСервис : %s\nЛогин : %s\nПароль : %s\nВремя до удаления пароля : %s", val.Service, val.Login, val.Password, val.Time)
	}
	return output
}
