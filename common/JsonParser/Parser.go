package JsonParser

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ReqMap map[string]string

func HTMLRequestBodyParser(req *http.Request) (ReqMap, error) {
	if nil == req.Body {
		return nil, nil
	}

	reqBodyBytes, err := ioutil.ReadAll(req.Body)
	if nil != err {
		log.Fatal(err)
		return nil, err
	}

	bodyString := string(reqBodyBytes)
	params := strings.Split(bodyString, "&")

	reqBody := make(map[string]string)

	if len(params) > 0 {
		for _, p := range params {
			param := strings.Split(p, "=")
			if len(param) > 1 {
				reqBody[param[0]] = param[1]
			}
		}
	} else {
		param := strings.Split(bodyString, "=")
		reqBody[param[0]] = param[1]
	}

	return reqBody, nil
}
