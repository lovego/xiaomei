package elastic

import (
	"net/http"
)

// 覆盖
func (es *ES) Put(path string, bodyData, data interface{}) error {
	resp, err := es.client.Put(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Check(http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	return resp.Json(data)
}

// 创建
func (es *ES) Create(path string, bodyData, data interface{}) error {
	resp, err := es.client.Put(es.Uri(path+`/_create`), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Check(http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	return resp.Json(data)
}

// 删除
func (es *ES) Delete(path string, data interface{}) error {
	return es.client.DeleteJson(es.Uri(path), nil, nil, data)
}

// 更新
func (es *ES) Update(path string, bodyData, data interface{}) error {
	return es.client.PostJson(es.Uri(path+`/_update`), nil, bodyData, data)
}

// Create if not Exist
func (es *ES) Ensure(path string, def interface{}) error {
	if ok, err := es.Exist(path); err != nil {
		return err
	} else if !ok {
		return es.Put(path, def, nil)
	}
	return nil
}

func (es *ES) Exist(path string) (bool, error) {
	resp, err := es.client.Head(es.Uri(path), nil, nil)
	if err != nil {
		return false, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, resp.CodeError()
	}
}
