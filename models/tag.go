package models



//标记
type Tag struct {
	Id int `redis:"id"`
	Tag string `redis:"tag"`
}



func (this *Tag) Load(id int)error{

	return nil
}