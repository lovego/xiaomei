package elastic

import (
	"bytes"
	"net/http"

	"github.com/lovego/xiaomei/utils/httputil"
)

func makeQueryDSL(must, filter []map[string]interface{}) []byte {
	boolQuery := make(map[string]interface{})
	if must != nil {
		boolQuery[`must`] = must
	}
	if filter != nil {
		boolQuery[`filter`] = filter
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{"bool": boolQuery},
	}
	buf, err := json.Marshal(query)
	if err != nil {
		PrintLog(`make query dsl error:`, err)
		return nil
	}
	return buf
}

/*
{
  "query": {
    "bool": {
      "must": [
        { "match": { "title":   "Search"        }},
        { "match": { "content": "Elasticsearch" }}
      ],
      "filter": [
        { "term":  { "status": "published" }},
        { "range": { "publish_date": { "gte": "2015-01-01" }}}
      ]
    }
  }
}
*/
