package HttpUtil

type ReqInfo struct {
	Protocol  string
	Method    string
	URL       string
	reqHeader map[string]string
	reqBody   []byte
	User      string
	Passwd    string
}

func NewRequest() *ReqInfo {
	req := &ReqInfo{
		Protocol:  "",
		Method:    "",
		URL:       "",
		reqHeader: make(map[string]string),
		reqBody:   make([]byte, 0, 0),
		User:      "",
		Passwd:    "",
	}

	return req
}

func (r *ReqInfo) SetProtocol(protocol string) {
	r.Protocol = protocol
}

func (r *ReqInfo) SetMethod(method string) {
	r.Method = method
}

func (r *ReqInfo) SetURL(url string) {
	r.URL = url
}

func (r *ReqInfo) AppendReqHeader(k, v string) {
	r.reqHeader[k] = v
}

func (r *ReqInfo) SetReqBody(body []byte) {
	r.reqBody = body
}

func (r *ReqInfo) SetUser(user string) {
	r.User = user
}

func (r *ReqInfo) SetPassword(pw string) {
	r.Passwd = pw
}
