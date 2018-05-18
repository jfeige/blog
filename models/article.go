package models


//文章
type Article struct {
	Id int `redis:"id"`
	Title string `redis:"title"`
	Content string `redis:"content"`
	Userid int  `redis:"userid"`
	Categoryid int `redis:"categoryid"`
	Tagid int `redis:"tagid"`
	Read_count int `redis:"read_count"`
	Comment_count int `redis:"comment_count"`
	Publish_time int `redis:"publish_time"`
	Publish_date int `redis:"publish_time"`
	Isshow int `redis:"isshow"`
}

//加载指定的文章
func (this *Article) Load(id int) error{

	return nil
}



//添加一篇文章
func (this *Article) Add() error{
	return nil
}