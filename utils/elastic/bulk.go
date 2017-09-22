package elastic

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data [][2]interface{}) error {
	if len(data) <= 0 {
		return nil
	}
	body, err := makeBulkCreate(data)
	if err != nil {
		return err
	}
	return es.BulkDo(path, body, `create`, data)
}

func (es *ES) BulkUpdate(path string, data [][2]interface{}) error {
	if len(data) <= 0 {
		return nil
	}
	return es.BulkDo(path, makeBulkUpdate(data), `update`, data)
}

func (es *ES) BulkDo(path string, body, typ string, data [][2]interface{}) error {
	result := BulkResult{}
	if err := es.client.PostJson(es.Uri(path+`/_bulk`), nil, body, &result); err != nil {
		return err
	}
	if !result.Errors {
		return nil
	}
	return bulkError{typ: typ, inputs: data, results: result.Items}
}

func (es *ES) BulkDelete(path string, data []interface{}) error {
	if len(data) <= 0 {
		return nil
	}

	result := BulkResult{}
	if err := es.client.PostJson(es.Uri(path+`/_bulk`), nil, makeBulkDelete(data), &result); err != nil {
		return err
	}
	if !result.Errors {
		return nil
	}
	return bulkDeleteError{inputs: data, results: result.Items}
}
