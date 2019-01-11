package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/freshautomations/telegram-moderator-bot/context"
	"os"
	"strconv"
	"time"
)

type UserData struct {
	Username string `json:"username"`
	UserID   int    `json:"id"`
	Name     string `json:"name"`
}

func Initialize(ctx *context.Context) {
	awscfg := aws.Config{
		Region: aws.String(ctx.Cfg.AWSRegion),
	}

	// Use IAM or environment variables credential
	if (os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "") ||
		(os.Getenv("AWS_ACCESS_KEY") != "" && os.Getenv("AWS_SECRET_KEY") != "") {
		awscfg.Credentials = credentials.NewEnvCredentials()
	}
	ctx.AWSSession = session.Must(session.NewSessionWithOptions(session.Options{Config: awscfg}))
	ctx.DDBSession = dynamodb.New(ctx.AWSSession)
	ctx.DBTableName = "tmb-" + ctx.Cfg.Environment + "-users"
}

func UpdateUserData(ctx *context.Context, User *UserData) (err error) {
	_, err = ctx.DDBSession.UpdateItem(&dynamodb.UpdateItemInput{
		//		ConditionExpression:         aws.String("attribute_not_exists #userid OR attribute_not_exists #name OR (attribute_exists #userid AND #userid <> :userid) OR (attribute_exists #name AND #name <> :name)"),
		ExpressionAttributeNames: map[string]*string{
			"#userid":   aws.String("id"),
			"#name":     aws.String("name"),
			"#lastseen": aws.String("lastseen"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userid": {
				N: aws.String(strconv.Itoa(User.UserID)),
			},
			":name": {
				S: aws.String(User.Name),
			},
			":lastseen": {
				N: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(User.Username),
			},
		},
		TableName:        aws.String(ctx.DBTableName),
		UpdateExpression: aws.String("SET #userid = :userid, #name = :name, #ttl = :ttl"),
	})
	return
}

func GetUserData(ctx *context.Context, username string) (*UserData, error) {
	result, err := ctx.DDBSession.GetItem(&dynamodb.GetItemInput{
		//		ConsistentRead: aws.Bool(true),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#username": aws.String("username"),
			"#userid":   aws.String("id"),
			"#name":     aws.String("name"),
		},
		ProjectionExpression: aws.String("#username, #userid, #name"),
		TableName:            aws.String(ctx.DBTableName),
	})
	if err != nil {
		return nil, err
	}

	output := UserData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &output)

	if output.UserID == 0 {
		return nil, err
	}

	return &output, err
}
