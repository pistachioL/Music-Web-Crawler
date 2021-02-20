package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func parseHtml(html string) {
	doc, err := goquery.NewDocument(html)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("a[class=pc_temp_songname]").Each(func(i int, selection *goquery.Selection) {
			selection.Attr("href")
			res:= selection.Text()
			fmt.Println(res)
	})



}
func main() {
	parseHtml("https://www.kugou.com/yy/html/rank.html")

}