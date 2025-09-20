package datasource

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/config"
)

type dynamoUserRepo struct {
	cli        *dynamodb.Client
	usersTable string
	idsTable   string
}

type userItem struct {
	UserID    int64  `dynamodbav:"userId"`
	Name      string `dynamodbav:"name"`
	Email     string `dynamodbav:"email"`
	Password  string `dynamodbav:"password"`
	CreatedAt int64  `dynamodbav:"createdAt"`
	UpdatedAt int64  `dynamodbav:"updatedAt"`
}

func NewDynamoUserRepository(ctx context.Context, cfg *config.Config) (port.UserRepository, error) {
	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(cfg.AWSRegion))
	if err != nil {
		return nil, err
	}
	return &dynamoUserRepo{cli: dynamodb.NewFromConfig(awsCfg), usersTable: cfg.UsersTableName, idsTable: cfg.IdsTableName}, nil
}

func (r *dynamoUserRepo) nextID(ctx context.Context, seq string) (int64, error) {
	out, err := r.cli.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.idsTable),
		Key: map[string]types.AttributeValue{
			"sequence": &types.AttributeValueMemberS{Value: seq},
		},
		UpdateExpression:          aws.String("SET #v = if_not_exists(#v, :zero) + :inc"),
		ExpressionAttributeNames:  map[string]string{"#v": "next"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":zero": &types.AttributeValueMemberN{Value: "0"}, ":inc": &types.AttributeValueMemberN{Value: "1"}},
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return 0, err
	}
	attr, ok := out.Attributes["next"].(*types.AttributeValueMemberN)
	if !ok {
		return 0, errors.New("invalid sequence response")
	}
	id, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *dynamoUserRepo) Create(ctx context.Context, u *domain.User) error {
	id, err := r.nextID(ctx, "user")
	if err != nil {
		return err
	}
	u.UserID = id
	item := userItem{
		UserID:    u.UserID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = r.cli.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.usersTable),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(userId)"),
	})
	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if errors.As(err, &cce) {
			return errors.New("user already exists")
		}
	}
	return err
}

func (r *dynamoUserRepo) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	res, err := r.cli.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.usersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, nil
	}
	var it userItem
	if err := attributevalue.UnmarshalMap(res.Item, &it); err != nil {
		return nil, err
	}
	u := &domain.User{UserID: it.UserID, Name: it.Name, Email: it.Email, Password: it.Password, CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt}
	return u, nil
}

func (r *dynamoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	res, err := r.cli.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.usersTable),
		IndexName:              aws.String("email_index"),
		KeyConditionExpression: aws.String("email = :e"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if res.Count == 0 || len(res.Items) == 0 {
		return nil, nil
	}
	var it userItem
	if err := attributevalue.UnmarshalMap(res.Items[0], &it); err != nil {
		return nil, err
	}
	u := &domain.User{UserID: it.UserID, Name: it.Name, Email: it.Email, Password: it.Password, CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt}
	return u, nil
}
