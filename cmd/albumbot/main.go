package main

import (
	"albumbot"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	const (
		tableKey = "TABLE_NAME"
		tokenKey = "DISCORD_TOKEN"
	)
	var (
		tableValue = flag.String(tableKey, "", "Name of DynamoDB Table.")
		tokenValue = flag.String(tokenKey, "", "Discord token.")
	)

	flag.Parse()
	godotenv.Load()

	m := map[string]string{
		tableKey: *tableValue,
		tokenKey: *tokenValue,
	}

	for k, v := range m {
		if err := overrideEnv(k, v); err != nil {
			panic(err)
		}
	}
}

// 与えられたキーの環境変数を与えられた文字列で上書きします。
// 元の環境変数と上書きする環境変数どちらも存在しない場合、errorを返します。
func overrideEnv(key, value string) error {
	if value != "" {
		os.Setenv(key, value)
		return nil
	} else if _, ok := os.LookupEnv(key); !ok {
		return fmt.Errorf("%sを指定してください。(-%s=<VALUE>)", key, key)
	}
	return nil
}

func main() {
	albumbot.New()
}
