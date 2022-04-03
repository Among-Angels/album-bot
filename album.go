package albumbot

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
	Title      string
	Urls       []string
	AlbumIndex int
}

// Albumsはアルバムデータが保存されたJsonファイル全体を表現します。
type Albums struct {
	Albums []Album
}

func getAlbumTitles(table string, c context.Context, client dynamodb.ScanAPIClient) (titles []string, e error) {
	var awsTable = aws.String(table)
	key := "Title, AlbumIndex"
	params := &dynamodb.ScanInput{
		TableName:            awsTable,
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
	sort.Slice(albums, func(i, j int) bool { return albums[i].AlbumIndex < albums[j].AlbumIndex })
	for _, al := range albums {
		titles = append(titles, al.Title)
	}
	return titles, nil
}

// GetAlbumTitlesはアルバム名のリストを返します。
func GetAlbumTitles(table string) (titles []string, e error) {
	return getAlbumTitles(table, context.TODO(), dbClient)
}

// GetAlbumUrlsは与えられたアルバム名の画像のURLのリストを返します。
func GetAlbumUrls(table string, title string) (urls []string, e error) {
	awsTable := aws.String(table)
	out, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: awsTable,
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
func GetAlbumPage(table, title string, start, count int) (urls []string, e error) {
	if start < 0 {
		e = fmt.Errorf("startは0以上の数値を指定してください")
		return nil, e
	}

	if count < 1 {
		e = fmt.Errorf("countは1以上の数値を指定してください")
		return nil, e
	}

	allUrls, err := GetAlbumUrls(table, title)
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

func titleExists(table, title string) (bool, error) {
	titles, err := GetAlbumTitles(table)
	if err != nil {
		return false, err
	}
	for _, t := range titles {
		if t == title {
			return true, nil
		}
	}
	return false, nil
}

// PostImageは与えられたアルバム名のUrl配列に与えられたUrlを追加します
func PostImage(table, title, url string) error {
	exists, err := titleExists(table, title)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("指定されたアルバム名は存在しません")
	}

	var awsTable = aws.String(table)
	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":url": &types.AttributeValueMemberSS{
				Value: []string{url},
			},
		},
		UpdateExpression: aws.String("ADD urls :url"),
		TableName:        awsTable,
	}
	_, err = dbClient.UpdateItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

// CreateAlbumは新しいアルバムをDynamoDB上に作成します
func CreateAlbum(table, title string) error {
	exists, err := titleExists(table, title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("すでに存在するアルバム名です。")
	}
	timestamp := fmt.Sprint(time.Now().Unix())
	var awsTable = aws.String(table)
	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"Title":      &types.AttributeValueMemberS{Value: title},
			"AlbumIndex": &types.AttributeValueMemberN{Value: timestamp},
		},
		TableName: awsTable,
	}
	_, err = dbClient.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAlbumは与えられたアルバム名をDynamoDB上から削除します
func DeleteAlbum(table, title string) error {
	exists, err := titleExists(table, title)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("存在しないアルバム名です。")
	}
	var awsTable = aws.String(table)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
		TableName: awsTable,
	}
	_, err = dbClient.DeleteItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAlbumUrlは与えられたアルバム名のUrlをDynamoDB上から削除します
func DeleteImage(table, title, url string) error {
	exists, err := titleExists(table, title)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("存在しないアルバム名です。")
	}
	var awsTable = aws.String(table)
	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"Title": &types.AttributeValueMemberS{Value: title},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":url": &types.AttributeValueMemberSS{
				Value: []string{url},
			},
		},
		UpdateExpression: aws.String("DELETE urls :url"),
		TableName:        awsTable,
	}
	_, err = dbClient.UpdateItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
