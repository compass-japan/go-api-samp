![test](https://github.com/ohorigenal/go-api-samp/workflows/test/badge.svg?branch=master&event=push)
![integration-test](https://github.com/ohorigenal/go-api-samp/workflows/integration-test/badge.svg?branch=master)

# Weather API
天気情報の登録、取得のAPI実装のサンプル。

フォーマットはJSON

言語：Go

FW：echo

他: MySQL,Docker,Kubernetes

docker-compose、minikube(local kubernetes)を利用しての実行は下記に記載してます

※ローカルのMySQLを利用する場合は、config.yamlの修正が必要となります

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
|location_id|int|○|"1"|
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

* 外部APIの東京の5日間の天候情報を整形したものを取得

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

## 実行(docker-compose)

```
# 事前準備
# dockerのinstall
https://docs.docker.com/get-docker/
```

リポジトリのルートディレクトリで実行

```
# buildと実行
$ docker-compose build --no-cache
$ docker-compose up

※イメージは１G少しありbuildのたびに増えていくので不要なイメージは削除も必要になると思います
$ docker rmi <image ID> (<image ID> ...) 
```

## 実行(minikube)

```
# 事前準備 
# dockerのinstall
https://docs.docker.com/get-docker/

# minikubeのinstall
https://kubernetes.io/ja/docs/tasks/tools/install-minikube/

# kubectlのinstall
https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/
```

```
# minikube起動
$ minikube start

# namespaceの作成
$ kubectl create namespace minikube

# contextの設定
$ kubectl config set-context $(kubectl config current-context) --namespace=minikube

# contextの確認
$ kubectl config get-contexts

# addonの追加(数分) https://kubernetes.io/ja/docs/tasks/access-application-cluster/ingress-minikube/
$ minikube addons enable ingress

# addonの追加を確認
$ kubectl get pods -n kube-system 

# minikubeのdockerを利用(minikube内にdocker imageを保存するため)
$ eval $(minikube docker-env)

# minikube環境を利用しているか確認(minikubeとなっていること)
$ docker info --format '{{json .Name}}'

# タグをつけてビルド
$ docker build  --no-cache -t go-api-samp/golang:v1.0 . 
$ docker images go-api-samp/golang:v1.0

# minikubeへのデプロイ
$ kubectl apply -k ./.k8s/overlay/minikube/

# ingressのADDRESS,HOSTを確認(host設定のため)
$ kubectl get ingress go-api-ingress

# ホストファイル/etc/hostsを編集 ※windowsではC:¥Windows¥System32¥drivers¥etc¥hosts
# macなどでvimで編集するなら
$ sudo vim /etc/hosts
例) 172.168.1.10 compass-j.com
確認したADDRESS HOSTを上記の例の形式で最終行に追加

--完了--
```

**アクセス例**
```
# URL
# docker-compose:localhost:8080
# minikube: compass-j.com

# /register
$ curl -D - -X POST -H 'Content-Type: application/json' -H 'Auth-Token: auth-token' -d '{"location_id":1, "date":"20200101", "weather":1, "comment":"good day"}' http://localhost:8080/register

# /get/:location_id/:date
$ curl -D - -X GET -H 'Auth-Token: auth-token' http://localhost:8080/get/1/20200101

# /get/apidata
$ curl -D - -X GET -H 'Auth-Token: auth-token' http://localhost:8080/get/apidata
```

**備考** 
```
# goの実行
# モジュールのダウンロード
$ go mod download

# 実行　※MySQLとの通信がうまくいかないと起動に失敗します
$ go run .

# テスト実行 ./...はカレント以下全てのtestファイル -vで詳しく
$ go test -v ./...
```
