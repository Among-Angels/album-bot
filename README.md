# albumbot

## bot の起動方法（ローカル版）

### 事前準備

1. .env をプロジェクト直下(albumbot/)に配置し、Discord トークンと DynamoDB のテーブル名を[リンク先](https://discordapp.com/channels/252122237761486849/788388972825018368/891657884570619934)のように書く
1. [AWS の認証情報を書く](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-files.html)

### 起動コマンド

```sh
go run ./cmd/albumbot/main.go
```

booted!が出たら成功

### Asset をダウンロードした場合

1. 実行前に AWS の認証情報を設定しておく

   - AWS CLI で設定する場合

   ```sh
   aws configure
   ```

   - 環境変数に設定する場合

   ```sh
   export AWS_ACCESS_KEY_ID=<アクセスキー>
   export AWS_SECRET_ACCESS_KEY=<シークレットアクセスキー>
   export AWS_DEFAULT_REGION=<リージョン>
   ```

   - もしくは、実行環境(EC2 インスタンスなど)に IAM ロールを与える

1. ダウンロードしたディレクトリに移動し、バイナリを実行する。

   ```sh
   ./main -TABLE_NAME=<DynamoDBのテーブル名> -DISCORD_TOKEN=<Discordトークン>
   ```

   or

   ```sh
   export TABLE_NAME=<DynamoDBのテーブル名>
   export DISCORD_TOKEN=<Discordトークン>
   ./main
   ```

   or

   ```sh
   echo TABLE_NAME=<DynamoDBのテーブル名> >> .env
   echo DISCORD_TOKEN=<Discordトークン> >> .env
   ./main
   ```

   booted!が出たら成功
