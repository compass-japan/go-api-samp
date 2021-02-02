# Weather API
天気情報の登録、取得、削除用のAPI実装のサンプル。

フォーマットはJSON

言語：Go

FW：echo

## IF仕様
#### 共通ヘッダ

**input**

|header|required|概要|
|---|---|---|
|Auth-Token|○|認証用(内部ネットワーク想定)、CASE_INSENSITIVE固定値("auth-token")|
|X-Request-ID| |ない場合作成 |

**output**

|header|type|概要|
|---|---|---|
|X-Request-ID|string| |

#### /register (POST)

* 天気情報の登録

**input**

|body key|type|required|概要|
|---|---|---|---|
|location_id|int| |"1"|
|date|string|○|"20200101" 八桁|
|weather|int|○|"1=sunny 2=cloudy 3=rainy 4=snowy"|
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
|location_id|int|○|1|
|date|string|○|"20200101" 八桁|

**output**

|key|type|概要|
|---|---|---|
|location|int|"新宿"|
|date|string|"20200101" 八桁|
|weather|string|"sunny cloudy rainy snow"|
|comment|string|一言コメント|

### /get/apidata (GET)

* 他のAPIから東京の今日の天候情報を取得する
https://www.metaweather.com/api/location/1118370/

**output**

|key| |type|概要|
|---|---|---|---|
|consolidated_weather| | |5日間の天気情報|
| |weather_state_name|string|天候名|
| |applicable_date|string|日|
| |wind_speed|float|風速|
| |air_pressure|float|大気圧|
| |humidity|int| 湿度 |
|title| |string| 都市 |
|timezone| |string| タイムゾーン |


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


