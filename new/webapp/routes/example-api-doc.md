# POST /example/path（API文档样例）

## 请求头说明
| Header       | 中文名称   | 必需 | 说明
| ------------ | --------   | ---- | -----------------
| Header1      | 请求头1    | 是   | Header1的说明

## Query参数说明
| 参数        | 中文名称     | 类型    | 必需 | 校验规则
| --------    | --------     | ------- | ---- | -------------------
| param1      | 参数1        | bool    | 否   | param1校验规则说明
| param2      | 参数2        | string  | 否   | param2校验规则说明

## 请求体说明 (application/x-www-form-urlencoded编码) （与json编码二选一，请删掉此备注）
| 参数        | 中文名称     | 类型    | 必需 | 校验规则
| --------    | --------     | ------- | ---- | -------------------
| param1      | 参数1        | bool    | 否   | param1校验规则说明
| param2      | 参数2        | string  | 否   | param2校验规则说明

## 请求体说明 (application/json编码) （与urlencoded编码二选一，请删掉此备注）
```
{
  "field1": "value1",    # 字段1的说明
  "field2": [            # 字段2的说明
    { 
      "field3": "value3",    # 字段3的说明
      "field4": "value4"     # 字段4的说明
    },
    ......
  ]
}
```

## 返回体说明 (application/json编码)
```
{
    "code": "ok",           # ok 表示成功，其他表示错误代码
    "message": "success"    # 与code对应的描述信息
}
```


## 请求示例
```
curl -XPOST 'https://example.com/exmaple/path?param1=v1&param2=v2' -H 'Header1: value' \
  -d'{
    "field1": "value1",
    "field2": [
      { "field3": "...", "field4": "..." },
      { "field3": "...", "field4": "..." }
    ]
  }'; echo
```

## 返回示例
```
{
    "code": "ok",
    "message": "success"
}
```
