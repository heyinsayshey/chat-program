package esQuery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"../common/HttpUtil"
	parser "../common/JsonParser"
)

var (
	QUERY_KEY_CREATE_DATE  = "$CREATE_DATE"
	QUERY_KEY_CREATE_TITLE = "$TITLE"
	QUERY_KEY_USERID       = "$USERID"
)

func Call(req *http.Request, filename, index string) (*http.Response, error) {
	queryBytes, err := buildQueryBody(req, filename)
	if nil != err {
		return nil, err
	}

	sendReq := HttpUtil.NewRequest()
	sendReq.SetMethod("POST")

	str := fmt.Sprintf("http://localhost:9200/%v", index)
	sendReq.SetURL(str)
	sendReq.AppendReqHeader("Content-Type", "application/json")
	sendReq.SetReqBody(queryBytes)
	sendReq.SetUser("root")
	sendReq.SetPassword("root")

	res, err := HttpUtil.SendRequest(sendReq)
	if nil != err {
		return nil, err
	}

	return res, nil
}

func buildQueryBody(req *http.Request, filename string) ([]byte, error) {
	reqBody, err := parser.HTMLRequestBodyParser(req)
	if nil != err {
		return nil, err
	}

	keyMap := make(map[string]string, 0)

	switch filename {
	case "createRoom.json":
		{
			title := reqBody["title"]
			unescapedTitle, _ := url.QueryUnescape(title)
			keyMap[QUERY_KEY_CREATE_TITLE] = unescapedTitle

			userId := "test1" // TODO : GET ID FROM SESSION !!!
			keyMap[QUERY_KEY_USERID] = userId

			now := time.Now()
			strTime := now.Format(time.RFC3339)
			idx := strings.Index(strTime, "+")
			strTime = strTime[:idx]
			strTime = strTime + ".000Z"
			keyMap[QUERY_KEY_CREATE_DATE] = strTime

			if title == "" || userId == "" {
				return nil, fmt.Errorf("Error : no matched parameter.")
			}
		}
	case "searchRoom.json":
		{
			userId := "test1" // TODO : GET ID FROM SESSION !!!
			keyMap[QUERY_KEY_USERID] = userId
		}
	case "deleteRoom.json":
		{

		}
	default:
		return nil, fmt.Errorf("Error : no matched filename.")
	}

	return getQueryFromFile(filename, keyMap)
}

func getQueryFromFile(filename string, keyMap map[string]string) ([]byte, error) {
	relativePath := fmt.Sprintf("./esQuery/%v", filename)
	absPath, _ := filepath.Abs(relativePath)
	content, e := ioutil.ReadFile(absPath)
	if nil != e {
		return nil, e
	}

	strContent := string(content)

	for k, v := range keyMap {
		strContent = strings.Replace(strContent, k, v, -1)
	}

	fmt.Println(strContent)

	return []byte(strContent), nil
}
