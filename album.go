package albumbot

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var table = aws.String(os.Getenv("TABLE_NAME"))
var dbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Create an Amazon DynamoDB client.
	dbClient = dynamodb.NewFromConfig(cfg)
}

// Albumはタイトルとそれに紐付けられた画像URLの集合です。
type Album struct {
	Title string
	Urls  []string
}

// Albumsはアルバムデータが保存されたJsonファイル全体を表現します。
type Albums struct {
	Albums []Album
}

func getAlbumTitles(c context.Context, client dynamodb.ScanAPIClient) (titles []string, e error) {
	key := "Title"
	params := &dynamodb.ScanInput{
		TableName:            table,
		ProjectionExpression: &key,
	}
	resp, err := client.Scan(c, params)
	if err != nil {
		fmt.Println("Got an error scanning the table:")
		fmt.Println(err.Error())
		return
	}
	albums := []Album{}
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &albums)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
	}
	for _, al := range albums {
		titles = append(titles, al.Title)
	}
	return titles, nil
}

// GetAlbumTitlesはアルバム名のリストを返します。
func GetAlbumTitles() (titles []string, e error) {
	return getAlbumTitles(context.TODO(), dbClient)
}

// GetAlbumUrlsは与えられたアルバム名の画像のURLのリストを返します。
func GetAlbumUrls(title string) (urls []string, e error) {
	out, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: table,
		Key: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
	})
	if err != nil {
		return nil, e
	}
	var album Album
	err = attributevalue.UnmarshalMap(out.Item, &album)
	if err != nil {
		return nil, e
	}
	return album.Urls, e
}

// GetAlbumPageは与えられたアルバム名と開始index,数量からURLのリストを返します
func GetAlbumPage(title string, start, count int) (urls []string, e error) {

	if start < 0 {
		e = fmt.Errorf("startは0以上の数値を指定してください")
		return nil, e
	}

	if count < 1 {
		e = fmt.Errorf("countは1以上の数値を指定してください")
		return nil, e
	}

	allUrls, err := GetAlbumUrls(title)
	if err != nil {
		return nil, err
	}

	if len(allUrls)-1 < start {
		e = fmt.Errorf("startはアルバム内のURL数より小さい値を指定してください")
		return nil, e
	}

	// count込で溢れた場合末尾まで返す
	if len(allUrls) < start+count {
		return allUrls[start:], nil
	}

	return allUrls[start : start+count], nil
}

// PostAlbumUrlは与えられたアルバム名のUrl配列に与えられたUrlを追加します
func PostAlbumUrl(title, url string) error {
	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":url": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: url},
				},
			},
		},
		UpdateExpression: aws.String("SET urls = list_append(urls, :url)"),
		TableName:        table,
	}
	_, err := dbClient.UpdateItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

// CreateAlbumは新しいアルバムをDynamoDB上に作成します
func CreateAlbum(title string) error {
	titles, err := GetAlbumTitles()
	if err != nil {
		return err
	}
	for _, t := range titles {
		if t == title {
			return fmt.Errorf("すでに存在するアルバム名です。")
		}
	}
	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
		TableName: table,
	}
	_, err = dbClient.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
