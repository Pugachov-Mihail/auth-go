package models

type User struct {
	Id       int64
	Login    string
	Email    string
	PassHash []byte
	SteamId  int64
}
