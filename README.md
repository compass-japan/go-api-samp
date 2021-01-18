# Weather API
天気情報の登録、取得、削除用のAPI実装のサンプル。

フォーマットはJSON

言語：Go

FW：echo

## IF仕様
#### 共通ヘッダ
|header|type|required|概要|
|---|---|---|---|
|Auth-Token|string|○|認証用(内部ネットワーク想定)、固定値|

#### /register (POST)

* 天気情報の登録

**input**

|body key|type|required|概要|
|---|---|---|---|
|date|string|○|"20200101" 八桁|
|weather|uint|○|"0=sunny 1=cloudy 2=rainy 3=snow"|
|location_id|string| |"1"|
|comment|string| |一言コメント|

**output**

|key|type|概要|
|---|---|---|
|message|string|"success"|



### /get/{location_id}/{date} (GET)

* 特定の日の天気情報の取得

※サンプルのためlocation_idは1固定

**input**

|path param|type|required|概要|
|---|---|---|---|
|location_id|string|○|"1"|
|date|string|○|"20200101" 八桁|

**output**

|key|type|概要|
|---|---|---|
|date|string|"20200101" 八桁|
|weather|string|"sunny cloudy rainy snow"|
|comment|string|一言コメント|

### /get/apidata (GET)

* 他のAPIからロンドンの今日の天候情報を取得する
https://www.metaweather.com/api/location/44418/

**output**

|key|type|概要|
|---|---|---|
|date|string|"20200101" 八桁|
|weather|string| api依存の天気情報 |


#### エラー共通レスポンス

|key|type|概要|
|---|---|---|
|message|string|"error message"|

## 実行例

**サーバ起動** 
```
$ go mod download
$ go run .
```

**アクセス**
```
# register
$ curl -D - -X POST -H 'Content-Type: application/json' -H 'Auth-Token: authtoken' -d '{"date":"20200101", "weather":"sunny", "comment":"good day"}' http://localhost:8080/register

# get
$ curl -D - -X GET -H 'Auth-Token: authtoken' http://localhost:8080/get/1/20200101
```


