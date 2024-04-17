package models

type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type Location struct {
	Longitude string `bson:"longitude" json:"longitude"`
	Latitude  string `bson:"latitude" json:"latitude"`
}
