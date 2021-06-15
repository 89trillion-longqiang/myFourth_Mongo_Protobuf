#1.创建礼品码

```
接口地址 
/giftCode/login 
```
## 请求方式
GET
## 请求示例
```
http://127.0.0.1:8080/giftCode/login?id=nccKwM9O
```
## 参数  说明

``` 
id 类型string 用户的id
```

```
成功示例 
{
    "condition": "success",
    "data": {
        "Diamond": "60",
        "Gold": "300",
        "Uid": "nccKwM9O"
    }
}
```

#2.验证礼品码接口

```
接口地址 
/giftCode/VerGiftCode 
```
## 请求方式
GET
## 请求示例
```
http://127.0.0.1:8080/giftCode/VerGiftCode?giftCode=nAyUzwZh&usr=nccKwM9O
```

## 参数  说明

``` 
giftCode 类型string 此字段为需要查询的礼包码
usr 类型string 用户的id
```

```
成功示例 
{
    "condition": "pass",
    "data": "GgQIARBkGgQIAhAUIgUIARCsAiIECAIQPCoECAIQUCoFCAEQkAM="
}

```
