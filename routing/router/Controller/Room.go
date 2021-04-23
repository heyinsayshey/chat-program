package Controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	parser "../../../common/JsonParser"
	r "../../renderer"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
)

type Room struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Create_date string        `json:"create_date"`
	Users       []interface{} `json:"users"`
}

func CreateRoom(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	//get Parameters
	if req.Body == nil {
		log.Fatal(fmt.Errorf("Error ::: empty request body."))
		return
	}

	reqBody, err := parser.HTMLRequestBodyParser(req)
	if nil != err {
		log.Fatal(err)
		return
	}

	encodedTitle := reqBody["title"]
	title, _ := url.QueryUnescape(encodedTitle)

	u := GetCurrentUser(req)
	userId := u.ID
	id := map[string]string{
		"id": userId,
	}
	ids := make([]interface{}, 0)
	ids = append(ids, id)

	now := time.Now()
	strTime := now.Format(time.RFC3339)
	idx := strings.Index(strTime, "+")
	strTime = strTime[:idx]
	strTime = strTime + ".000Z"

	var buf bytes.Buffer
	buf.WriteString(`{"create_date" : "`)
	buf.WriteString(strTime)
	buf.WriteString(`", "title" : "`)
	buf.WriteString(title)
	buf.WriteString(`", "users" : [{ "id" : "`)
	buf.WriteString(userId)
	buf.WriteString(`" }]}`)
	strQuery := buf.String()
	fmt.Println(strQuery)

	// create new client
	ctx := context.Background()

	esCli, err := elastic.NewClient()
	if nil != err {
		log.Fatal(err)
		return
	}

	createRes, err := esCli.Index().
		Index("rooms").
		BodyString(strQuery).
		Do(ctx)
	if nil != err {
		log.Fatal(err)
		return
	}

	fmt.Println(createRes)

	var room Room
	if strings.Compare(createRes.Result, "created") == 0 {
		room = Room{
			ID:          createRes.Id,
			Create_date: strTime,
			Title:       title,
			Users:       ids,
		}
	}
	_, err = esCli.Flush("rooms").Do(ctx)
	if nil != err {
		log.Fatal(err)
	}

	ren := r.GetInstance()
	ren.JSON(w, http.StatusOK, room)
}

func SearchRooms(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	// create new client
	ctx := context.Background()

	esCli, err := elastic.NewClient()
	if nil != err {
		log.Fatal(err)
		return
	}

	//build Query
	u := GetCurrentUser(req)

	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewTermQuery("users.id", u.ID))
	q := elastic.NewNestedQuery("users", bq)
	src, err := q.Source()
	if err != nil {
		log.Fatal(err)
		return
	}

	data, _ := json.Marshal(src)
	fmt.Println("[query] ::: ", string(data))

	// send query
	searchRes, err := esCli.Search().
		Index("rooms").
		Query(q).
		Sort("create_date", true).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if nil != err {
		log.Fatal(err)
		return
	}

	fmt.Printf("Found total of %d rooms\n", searchRes.TotalHits())
	var arrRoom []Room
	if searchRes.TotalHits() > 0 {
		for _, hit := range searchRes.Hits.Hits {

			var r map[string]interface{}
			err := json.Unmarshal(hit.Source, &r)
			if nil != err {
				log.Fatal(err)
				return
			}

			newRoom := Room{
				ID:          hit.Id,
				Create_date: r["create_date"].(string),
				Title:       r["title"].(string),
				Users:       r["users"].([]interface{}),
			}

			//fmt.Println(newRoom)
			arrRoom = append(arrRoom, newRoom)
		}
	}
	_, err = esCli.Flush("rooms").Do(ctx)
	if nil != err {
		log.Fatal(err)
	}

	ren := r.GetInstance()
	ren.JSON(w, http.StatusOK, arrRoom)
}

func DeleteRoom(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

}
