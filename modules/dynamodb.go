package main

import (
	"strconv"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type BackendDynamoDB struct {
	config     *stns.Config
	db         *dynamodb.DynamoDB
	userTable  string
	groupTable string
}

func NewBackendDynamodb(c *stns.Config) (model.Backend, error) {
	readCapacityUnits := int64(5)
	writeCapacityUnits := int64(5)
	userTable := "stns_users"
	groupTable := "stns_groups"

	if c.Modules["dynamodb"] != nil {
		if c.Modules["dynamodb"].(map[string]interface{})["read_capacity_units"] != nil {
			readCapacityUnits = c.Modules["dynamodb"].(map[string]interface{})["read_capacity_units"].(int64)
		}

		if c.Modules["dynamodb"].(map[string]interface{})["write_capacity_units"] != nil {
			writeCapacityUnits = c.Modules["dynamodb"].(map[string]interface{})["write_capacity_units"].(int64)
		}

		if c.Modules["dynamodb"].(map[string]interface{})["user_table_name"] != nil {
			userTable = c.Modules["dynamodb"].(map[string]interface{})["user_table_nane"].(string)
		}

		if c.Modules["dynamodb"].(map[string]interface{})["group_table_name"] != nil {
			groupTable = c.Modules["dynamodb"].(map[string]interface{})["group_table_nane"].(string)
		}
	}

	db := dynamodb.New(session.New())
	tables, err := db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}

CREATE_TABLE:
	for _, tname := range []string{userTable, groupTable} {
		for _, n := range tables.TableNames {
			if *n == tname {
				continue CREATE_TABLE
			}
		}

		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("N"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(readCapacityUnits),
				WriteCapacityUnits: aws.Int64(writeCapacityUnits),
			},
			TableName: aws.String(tname),
		}

		_, err := db.CreateTable(input)
		if err != nil {
			return nil, err
		}

	}
	b := BackendDynamoDB{
		config:     c,
		db:         db,
		userTable:  userTable,
		groupTable: groupTable,
	}

	if c.Modules["dynamodb"] != nil && c.Modules["dynamodb"].(map[string]interface{})["sync"] != nil && c.Modules["dynamodb"].(map[string]interface{})["sync"].(bool) {
		err := syncConfig(b, c)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b BackendDynamoDB) findAll(table, resource string) ([]map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}

	result, err := b.db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}

	return result.Items, nil
}

func (b BackendDynamoDB) findByName(table, resource, value string) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {
				S: aws.String(value),
			},
		},
		FilterExpression: aws.String("#name = :value"),
		TableName:        aws.String(table),
	}

	result, err := b.db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}

	return result.Items[0], nil
}

func (b BackendDynamoDB) findByID(table, resource, id string) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String(table),
	}

	result, err := b.db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}
	return result.Item, nil
}

func (b BackendDynamoDB) FindUserByName(name string) (map[string]model.UserGroup, error) {
	item, err := b.findByName(b.userTable, "user", name)
	if err != nil {
		return nil, err
	}

	users := model.Users{}
	user := new(model.User)

	if err := dynamodbattribute.UnmarshalMap(item, user); err != nil {
		return nil, err
	}
	users[user.GetName()] = user
	return users.ToUserGroup(), nil
}

func (b BackendDynamoDB) FindUserByID(id int) (map[string]model.UserGroup, error) {
	item, err := b.findByID(b.userTable, "user", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	users := model.Users{}
	user := new(model.User)
	if err := dynamodbattribute.UnmarshalMap(item, user); err != nil {
		return nil, err
	}
	users[user.GetName()] = user
	return users.ToUserGroup(), nil
}

func (b BackendDynamoDB) Users() (map[string]model.UserGroup, error) {
	items, err := b.findAll(b.userTable, "user")
	if err != nil {
		return nil, err
	}

	users := model.Users{}
	user := new(model.User)
	for _, item := range items {
		if err := dynamodbattribute.UnmarshalMap(item, user); err != nil {
			return nil, err
		}
		users[user.GetName()] = user
	}
	return users.ToUserGroup(), nil
}

func (b BackendDynamoDB) FindGroupByName(name string) (map[string]model.UserGroup, error) {
	item, err := b.findByName(b.groupTable, "group", name)
	if err != nil {
		return nil, err
	}

	groups := model.Groups{}
	group := new(model.Group)

	if err := dynamodbattribute.UnmarshalMap(item, group); err != nil {
		return nil, err
	}
	groups[group.GetName()] = group
	return groups.ToUserGroup(), nil
}

func (b BackendDynamoDB) FindGroupByID(id int) (map[string]model.UserGroup, error) {
	item, err := b.findByID(b.groupTable, "group", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	groups := model.Groups{}
	group := new(model.Group)
	if err := dynamodbattribute.UnmarshalMap(item, group); err != nil {
		return nil, err
	}
	groups[group.GetName()] = group
	return groups.ToUserGroup(), nil
}

func (b BackendDynamoDB) Groups() (map[string]model.UserGroup, error) {
	items, err := b.findAll(b.groupTable, "group")
	if err != nil {
		return nil, err
	}

	groups := model.Groups{}
	group := new(model.Group)
	for _, item := range items {
		if err := dynamodbattribute.UnmarshalMap(item, group); err != nil {
			return nil, err
		}
		groups[group.GetName()] = group
	}
	return groups.ToUserGroup(), nil

}

func (b BackendDynamoDB) highlowID(table string, high bool) int {
	ret := 0
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#ID": aws.String("id"),
		},
		ProjectionExpression: aws.String("#ID"),
		TableName:            aws.String(table),
	}

	result, err := b.db.Scan(input)
	if err != nil {
		return ret
	}
	if len(result.Items) == 0 {
		return ret
	}
	for _, item := range result.Items {
		id, err := strconv.Atoi(*item["id"].N)
		if err != nil {
			return 0
		}
		if ret == 0 || (high && id > ret) || (!high && id < ret) {
			ret = id
		}
	}
	return ret
}

func (b BackendDynamoDB) HighestUserID() int {
	return b.highlowID(b.userTable, true)
}

func (b BackendDynamoDB) LowestUserID() int {
	return b.highlowID(b.userTable, false)
}

func (b BackendDynamoDB) HighestGroupID() int {
	return b.highlowID(b.groupTable, true)
}

func (b BackendDynamoDB) LowestGroupID() int {
	return b.highlowID(b.groupTable, false)
}

func (b BackendDynamoDB) CreateUser(v model.UserGroup) error {
	return b.update(b.userTable, v)
}

func (b BackendDynamoDB) CreateGroup(v model.UserGroup) error {
	return b.update(b.groupTable, v)
}

func (b BackendDynamoDB) update(table string, v model.UserGroup) error {
	av, err := dynamodbattribute.MarshalMap(v)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	_, err = b.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (b BackendDynamoDB) DeleteUser(id int) error {
	return b.delete(b.userTable, id)
}

func (b BackendDynamoDB) DeleteGroup(id int) error {
	return b.delete(b.groupTable, id)
}

func (b BackendDynamoDB) delete(table string, id int) error {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	}
	_, err := b.db.DeleteItem(params)
	if err != nil {
		return err
	}
	return nil
}

func (b BackendDynamoDB) UpdateUser(v model.UserGroup) error {
	return b.update(b.userTable, v)
}

func (b BackendDynamoDB) UpdateGroup(v model.UserGroup) error {
	return b.update(b.groupTable, v)
}
