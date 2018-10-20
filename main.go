package main

import (
	"encoding/json"
	"fmt"
	"inputer/inputInterfaces"
	"io/ioutil"
	"net/http"
)

type TaskPageJson struct {
	Data          []TaskDataJson
	Next_page_url string
}

type TaskDataJson struct {
	Category_id    int
	Daily_num      int
	Site           Site
	Input_category json.RawMessage
}

type Site struct {
	Input_url    string
	Input_module string
}

type Task struct {
	CategoryId    int
	DailyNum      int
	Site          Site
	InputCategory json.RawMessage
}

type Message struct {
	Status  string
	Code    int
	Message string
}

var taskList []Task

var serverUrl = "http://writer.localhost"

func main() {
	// 获取所有任务
	for {
		next := taskget(serverUrl + "/api/inputer/combination/task")
		if next == "" {
			break
		}
	}

	// 遍历任务发布文章
	for _, task := range taskList {
		category := inputInterfaces.DestoonCategoryParse(task.InputCategory) // 解析任务发布分类信息

		// 调用发布接口
		inputer := &inputInterfaces.Destoon{
			Url:           task.Site.Input_url,
			InputCategory: category,
		}

		articles := articleget(task.CategoryId, task.DailyNum) // 获取任务文章

		inputer.InputAll(&articles) // 发布文章
	}

}

// 获取任务列表
func taskget(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var page TaskPageJson
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &page)
	if err != nil {
		panic(err)
	}

	// 加入任务列表
	for _, item := range page.Data {
		task := Task{
			CategoryId:    item.Category_id,
			DailyNum:      item.Daily_num,
			Site:          item.Site,
			InputCategory: item.Input_category,
		}
		taskList = append(taskList, task)
	}

	return page.Next_page_url
}

// 获取分类文章列表
func articleget(catid, num int) []inputInterfaces.Article {
	url := fmt.Sprintf(serverUrl+"/api/inputer/combination/article/limit?catid=%d&num=%d&status=0", catid, num)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var articles []inputInterfaces.Article
	json.Unmarshal(body, &articles)
	return articles
}
