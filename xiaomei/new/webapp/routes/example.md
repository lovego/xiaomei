# POST /example/path（API文档样例）

## 请求头说明
| Header名称   | Header值   | 必需 | 说明
| ------------ | --------   | ---- | -----------------
| Cookie       | token=XXXX | 是   | 改请求头的说明

## Query参数说明
| 参数        | 中文名称     | 类型    | 必需 | 校验规则
| --------    | --------     | ------- | ---- | -------------------
| param1      | 参数1        | bool    | 否   | param1校验规则说明
| param2      | 参数2        | string  | 否   | param2校验规则说明

## 请求体说明 (application/x-www-form-urlencoded编码)
| 参数        | 中文名称     | 类型    | 必需 | 校验规则
| --------    | --------     | ------- | ---- | -------------------
| param1      | 参数1        | bool    | 否   | param1校验规则说明
| param2      | 参数2        | string  | 否   | param2校验规则说明

## 请求体说明 (application/json编码)
```
{
  "field1": "value1",    # 字段1的说明
  "field2": [            # 字段2的说明
    { 
      "field3": "value3",    # 字段3的说明
      "field4": "value4",    # 字段4的说明
    },
    ......
  ]
}
```

## 返回体说明
```
{
    "code": "ok",           # ok 表示成功，其他表示错误代码
    "message": "success"    # 用户友好的描述性信息，可以直接显示给用户
}
```


## 请求示例
```
curl -XPOST 'https://goods-stocks.qa.hztl3.com/company-parts/alliances' \
  -H 'Cookie: token=MTUxNzMwNjc5OHxleUpWYzJWeVNXUWlPakV3TURJMWZRbz18rPAiX1orjeaL4RfwRkmV5GZobu2jBR7Vbf19obxORI8=' \
  -d'{
		"swPartIds": [1001, 1002, 1003, 1004],
		"allianceIds": [21, 22]
	}'; echo
```

## 返回示例
```
{
    "code": "ok",
    "message": "success"
}
```
