package HttpUtil

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

func SendRequest(req *ReqInfo) (*http.Response, error) {
	r, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.reqBody))
	if nil != err {
		return nil, err
	}

	r.SetBasicAuth(req.User, req.Passwd)

	for k, v := range req.reqHeader {
		r.Header.Add(k, v)
	}

	// set tls enable
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	r.Close = true

	return cli.Do(r)
}
