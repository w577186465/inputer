package inputInterfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DestoonCategory struct {
	Moduleid int
	Catid    int
}

type Destoon struct {
	Url           string
	UserName      string
	Password      string
	InputCategory DestoonCategory
	SuccessIds    []int
}

var serverUrl = "http://writer.localhost"

func DestoonCategoryParse(raw json.RawMessage) DestoonCategory {
	jsonByte, err := raw.MarshalJSON()
	if err != nil {
	}
	var category DestoonCategory
	json.Unmarshal(jsonByte, &category)

	return category
}

func (info *Destoon) Input(article Article) {
	midValue := fmt.Sprintf("%d", info.InputCategory.Moduleid)
	catidValue := fmt.Sprintf("%d", info.InputCategory.Catid)

	values := url.Values{
		"mid":           {midValue},
		"post[title]":   {article.Title},
		"post[catid]":   {catidValue},
		"post[content]": {article.Content},
	}
	formData := values.Encode()

	resp, err := http.Post(info.Url, "application/x-www-form-urlencoded", strings.NewReader(formData))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析返回消息
	var message Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		// 发布出错
		return
	}

	if message.Status == "success" {
		info.SuccessIds = append(info.SuccessIds, article.Id)
	}

}

func (info *Destoon) InputAll(articles *[]Article) {
	for _, article := range *articles {
		info.Input(article)
	}

	info.StatusUpdate()
}

func (info *Destoon) StatusUpdate() {
	submitUrl := serverUrl + "/api/inputer/combination/article/inputed"
	formData := map[string]interface{}{
		"ids": info.SuccessIds,
	}

	bytesData, err := json.Marshal(formData)
	resp, err := http.Post(submitUrl, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
