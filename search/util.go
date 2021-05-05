package search

import (
	"context"
	"fmt"
	"githubLogin/model"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
	"strconv"
)
type Subject struct {
	Keyword string `json:"keyword"`
}
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

//存储搜索的歌曲结果
func saveSearchRes(searchResArr []model.DetailReq) {
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

//存储关键词 统计搜索次数 completion suggester
func SaveSearchKey() {
	client := connES()
	const mapping = `
	{
		"mappings": {
			"type_name": {    
				"properties": {
					"keyword" : {
						"type": "completion",
          				"analyzer": "standard"
					}
				}
			}	
		}
	}`
	//创建关键词索引
	//createIdx, err  := client.CreateIndex("keyword").BodyString(mapping).Do(context.Background())
	//if err != nil {
	//	fmt.Println("创建关键词索引错误", err)
	//	return
	//}
	//fmt.Println("索引创建成功", createIdx)

	//写入建议的数据（关键词）
	subject := Subject{
		Keyword: "邓紫棋",
	}
	subject1 := Subject{
		Keyword: "周杰伦",
	}

	doc, err := client.Index().
		Index("keyword").
		Type("keyword").
		Id("1").
		BodyJson(subject).
		Do(context.Background())

		client.Index().
		Index("keyword").
		Type("keyword").
		Id("2").
		BodyJson(subject1).
		Do(context.Background())

	if err != nil {
		fmt.Println("写入关键词失败：",err)
		return
	}
	fmt.Printf("Indexed with type=%s\n",doc.Type)

	//获取关键词
	res, err := client.Get().Index("keyword").Id("1").Do(context.Background())
	if(err != nil) {
		fmt.Println("获取关键词错误：",err)
		return
	}
	fmt.Println("关键词", res)
	if res.Found {
		fmt.Printf("Got document %v (version=%d, index=%s, type=%s)\n",
			res.Id, res.Version, res.Index, res.Type)
		//err := json.Unmarshal(res, &subject)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(subject.Id, subject.Keyword)
	}
	searchWithSuggest()
}


func searchWithSuggest() {
	client := connES()
	termQuery := elastic.NewTermQuery("keyword", "邓紫棋")
	searchResult, err := client.Search().
		Index("keyword").
		Query(termQuery).
		//Sort("id", true). // 按id升序排序
		//From(0).Size(10). // 拿前10个结果
		Pretty(true).
		Do(context.Background()) // 执行
	if err != nil {
		fmt.Println("搜索失败", err)
	}

	fmt.Printf("Found %d subjects\n", searchResult.TookInMillis)
	if searchResult.TookInMillis > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(Subject{})) {
			if t, ok := item.(Subject); ok {
				fmt.Println("Found: Subject(title=%s)\n", t.Keyword)
			}

		}

	} else {
		fmt.Println("Not found!")
	}
}


//根据关键词从es中查询结果
func getSearchResult(keyword string) {
	client := connES()
	searchRes,err := client.Get().Index("search").Type("search").Id(keyword).Do(context.Background())
	if err != nil {
		fmt.Println("ES获取搜索数据失败：", err)
	}
	fmt.Println("搜索结果：", searchRes)
}

func searchSuggest() {

}


//func getAllRes(res *elastic.SearchResult, err error) {
//	if err != nil {
//		print(err.Error())
//		return
//	}
//	var typ DetailReq
//	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
//		t := item.(DetailReq)
//		fmt.Printf("%#v\n", t)
//	}
//}
//

