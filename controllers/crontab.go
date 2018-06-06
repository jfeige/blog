package controllers

import (
	"blog/models"
	"sync"
	"time"
	"fmt"
)


var(
	readData map[int]int
	lock sync.RWMutex
)

func AddReadCnt(a_id int){
	if models.ReadChan == nil{
		models.ReadChan = make(chan int,1000)
	}
	models.ReadChan <- a_id
	fmt.Println("收到一条浏览请求:",a_id)
}
/**
	文章浏览数入库
 */
func ProcessReadData(){
	if readData == nil{
		readData = make(map[int]int)
	}
	if models.ReadChan == nil{
		models.ReadChan = make(chan int,1000)
	}
	for{
		select{
		case a_id :=<- models.ReadChan:
			//执行入库更新
			go addRead(a_id)
		case <- time.After(10*time.Second):
			//执行入库
			go updateReadCnt()
		}
	}
}

func addRead(a_id int){
	lock.Lock()
	defer lock.Unlock()
	v, ok := readData[a_id]
	if !ok {
		readData[a_id] = 1
	} else {
		readData[a_id] = v + 1
	}
}

func updateReadCnt(){
	lock.Lock()
	defer lock.Unlock()
	for a_id,cnt := range readData{
		//执行入库
		go models.UpdateReadCnt(a_id,cnt)
	}

	readData = make(map[int]int)
}