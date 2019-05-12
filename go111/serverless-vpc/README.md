# Serverless VPC Access Sample
## Serverless VCP Accessに対応しているサービス
* Google App Engine Standard Environment
  - PHP 5.5, Go 1.9を除く
* Cloud Functions


## Serverless VPC Accessが有効そうな場面
* 全文検索エンジンをGKEなどを使ってデプロイしている
  - わざわざサービスアカウントなどを使ってリクエスト認可のためのトークンを生成したりしてた
* Memorystoreを利用する
  - そもそもGAE/SEなどからは利用できなかった


## 注意点
* まだBeta版のため、Serverless VPC Access Connectorは `us-central1` した対応してない
  - GCFやGAEのリージョンは、Serverless VPC Access Connectorと同じリージョンにする必要あり
  - **App Engineを東京リージョンなどで設定してしまっているプロジェクトでは利用不可**
* ConnectorのIP範囲は、VPCで未使用の範囲を指定する必要がある
* GCEやGKEからGAEへのアクセスは内部IPを使うことはできない
* Connectorにネットワークタグをつけることはできない
