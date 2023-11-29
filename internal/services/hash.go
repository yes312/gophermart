package services

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(login, password string) string {
	// TODO получение хеша нужно усложнить. добавить соль, может еще как-то

	// #ВОПРОСМЕНТОРУ что делать с солью? хранить в yaml файле, считывать в конфиг перед запуском сервера,
	// потом передавать в хендлере?
	// или сделать метод, в котором будем хранить соль и его добавить в хендлер? наверное этот вариант лучше

	hash := sha256.Sum256([]byte(login + password))
	return hex.EncodeToString(hash[:])

}
