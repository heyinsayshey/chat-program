package Controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	parser "../../../common/JsonParser"
	"../../../common/RDBMS"
	r "../../renderer"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
)

type Message struct {
	ID        string `json:"id"`
	Room_id   string `json:"room_id"`
	Content   string `json:"content"`
	Send_date string `json:"send_date"`
	User      User   `json:"user"`
}

const LIMIT_MAX int = 10

func (m *Message) CreateMessage() error {

	//build query
	room_id := m.Room_id
	content := m.Content
	send_date := m.Send_date
	user := m.User.ID

	strQuery := fmt.Sprintf(`{
		"room_id" : "%s",
		"content" : "%s",
		"send_date" : "%s",
		"users" : "%s"
	}
	`, room_id, content, send_date, user)
	fmt.Println(strQuery)

	// create new client
	ctx := context.Background()

	esCli, err := elastic.NewClient()
	if nil != err {
		return err
	}

	_, err = esCli.Index().
		Index("messages").
		BodyString(strQuery).
		Do(ctx)
	if nil != err {
		return err
	}

	return nil
}

func SearchMessages(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	if req.Body == nil {
		log.Fatal(fmt.Errorf("Error ::: empty request body."))
		return
	}

	reqBody, err := parser.HTMLRequestBodyParser(req)
	if nil != err {
		log.Fatal(err)
		return
	}

	limit := LIMIT_MAX
	encodedRoomId := reqBody["id"]
	roomId, _ := url.QueryUnescape(encodedRoomId)

	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewMatchQuery("room_id", roomId))

	ctx := context.Background()

	esCli, err := elastic.NewClient()
	if nil != err {
		log.Fatal(err)
		return
	}

	searchRes, err := esCli.Search().
		Index("messages").
		Query(bq).
		Sort("send_date", false).
		From(0).Size(limit).
		Pretty(true).
		Do(ctx)
	if nil != err {
		log.Fatal(err)
		return
	}

	var arrMsgs []Message
	if searchRes.TotalHits() > 0 {
		for _, hit := range searchRes.Hits.Hits {

			var m map[string]interface{}
			err := json.Unmarshal(hit.Source, &m)
			if nil != err {
				log.Fatal(err)
				return
			}

			// get pic url from rdbms
			userid := m["users"].(string)
			sql := fmt.Sprintf("select thumbnail_path from tbl_member where id = '%v';", userid)
			val, err := RDBMS.MySQLSelect(sql)
			if nil != err {
				log.Fatal(err)
				return
			}
			//

			newMessage := Message{
				ID:        hit.Id,
				Room_id:   roomId,
				Content:   m["content"].(string),
				Send_date: m["send_date"].(string),
				User: User{
					ID:             m["users"].(string),
					Thumbnail_path: val,
					Timestamp:      time.Now(),
				},
			}

			fmt.Println(newMessage)
			arrMsgs = append(arrMsgs, newMessage)
		}
	}
	_, err = esCli.Flush("messages").Do(ctx)
	if nil != err {
		log.Fatal(err)
	}

	ren := r.GetInstance()
	ren.JSON(w, http.StatusOK, arrMsgs)
}
