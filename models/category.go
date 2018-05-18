package models


//类别
type Category struct {
	Id int `redis:"id"`
	Name string `redis:"name"`
	Index string `redis:"index"`
}