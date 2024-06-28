package knight

type Response struct {
	to   string
	data []byte
}

func New() *Response {
	return &Response{}
}

func (res *Response) To(to string) {
	res.to = to
}

func (res *Response) Write(data []byte) {
	if res == nil {
		res = New()
	}
	res.data = data
}
