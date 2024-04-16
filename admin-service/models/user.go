package models

type Admin struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}
