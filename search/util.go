package search

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
	"strconv"
)

func connES() *elastic.Client{
	// 创建ES client用于后续操作ES
	client, err := elastic.NewClient(
		// 设置ES服务地址，支持多个地址
		elastic.SetURL("http://127.0.0.1:9200"))
		// 设置基于http base auth验证的账号和密码
		//elastic.SetBasicAuth("user", "secret"))
	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
	} else {
		fmt.Println("连接成功")
	}
	return client
}

func saveSearchRes(searchResArr []DetailReq) {
	client := connES()
	n := 0
	bulkRequest := client.Bulk()
	for searchRes := range searchResArr {
		n++
		req := elastic.NewBulkIndexRequest().Index("search").Type("search").Id(strconv.Itoa(n)).Doc(searchRes)
		bulkRequest = bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		fmt.Println("存入ES错误:",err)
		return
	}
	fmt.Println("bulkResponse:",bulkResponse)
}

//存储关键词 统计搜索次数
func saveSearchKey(key string) {

}

func getSearchResult() {
	client := connES()
	searchRes,err := client.Search("search").Type("search").Do(context.Background())
	if err != nil {
		fmt.Println("ES获取搜索数据失败：", err)
	}
	fmt.Println("搜索结果：")
	printRes(searchRes, err)
}


//打印查询到结果
func printRes(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ DetailReq
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(DetailReq)
		fmt.Printf("%#v\n", t)
	}
}
