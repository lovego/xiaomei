# POST /example/path（API文档样例）

## 请求头说明
| 请求头       | 中文名称   | 类型   | 必需 | 说明                                  |
| ------------ | --------   | ------ | ---- | ------------------------------------- |
| Timestamp    | 时间戳     | int64  | 是   | Unix时间戳，与服务器时间误差1分钟以内 |
| Sign         | 签名       | string | 是   | 64位HEX编码的SHA256签名，函数示意：HEX(SHA256("timestamp,secret"))，签名内容是以","分隔的时间戳和密钥 |

## Query参数说明
| 参数         | 中文名称   | 类型   | 必需 | 校验规则            |
| --------     | --------   | ------ | ---- | ------------------- |
| param1       | 参数1      | bool   | 否   | param1校验规则说明  |
| param2       | 参数2      | string | 否   | param2校验规则说明  |

## 请求体说明 (application/json编码) 
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
