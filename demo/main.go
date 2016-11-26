package main

func layoutData(layout string, data interface{}, req *Request, res *Response) interface{} {
	if strings.HasPrefix(layout, `layout/`) {
		return struct {
			Data interface{}
		}{data}
	}
	return data
}
