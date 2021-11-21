# albumbot

## botの起動方法（ローカル版）

### 事前準備

1. .envをプロジェクト直下(albumbot/)に配置し、[リンク先](https://discordapp.com/channels/252122237761486849/788388972825018368/891657884570619934)のように書く
1. [AWSの認証情報を書く](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-files.html)

### 起動コマンド

```sh
go run ./cmd/albumbot/main.go  
```

booted!が出たら成功
