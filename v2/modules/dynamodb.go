package main

import (
	"context"
	"strconv"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BackendDynamoDB struct {
	config     *stns.Config
	db         *dynamodb.Client
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
			userTable = c.Modules["dynamodb"].(map[string]interface{})["user_table_name"].(string)
		}

		if c.Modules["dynamodb"].(map[string]interface{})["group_table_name"] != nil {
			groupTable = c.Modules["dynamodb"].(map[string]interface{})["group_table_name"].(string)
		}
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	db := dynamodb.NewFromConfig(cfg)
	tables, err := db.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}

CREATE_TABLE:
	for _, tname := range []string{userTable, groupTable} {
		for _, n := range tables.TableNames {
			if n == tname {
				continue CREATE_TABLE
			}
		}

		tnameValue := tname
		idAttr := "id"
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: &idAttr,
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: &idAttr,
					KeyType:       types.KeyTypeHash,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  &readCapacityUnits,
				WriteCapacityUnits: &writeCapacityUnits,
			},
			TableName: &tnameValue,
		}

		_, err := db.CreateTable(context.Background(), input)
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

func (b BackendDynamoDB) findAll(table, resource string) ([]map[string]types.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		TableName: &table,
	}

	result, err := b.db.Scan(context.Background(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}

	return result.Items, nil
}

func (b BackendDynamoDB) findByName(table, resource, value string) (map[string]types.AttributeValue, error) {
	nameAttr := "name"
	filterExpr := "#name = :value"
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]string{
			"#name": nameAttr,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":value": &types.AttributeValueMemberS{
				Value: value,
			},
		},
		FilterExpression: &filterExpr,
		TableName:        &table,
	}

	result, err := b.db.Scan(context.Background(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}

	return result.Items[0], nil
}

func (b BackendDynamoDB) findByID(table, resource, id string) (map[string]types.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{
				Value: id,
			},
		},
		TableName: &table,
	}

	result, err := b.db.GetItem(context.Background(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, model.NewNotFoundError(resource, nil)
	}
	return result.Item, nil
}

func (b BackendDynamoDB) unmarshalUserGroup(userGroup model.UserGroup, userGroups map[string]model.UserGroup, item map[string]types.AttributeValue) (map[string]model.UserGroup, error) {
	if err := attributevalue.UnmarshalMap(item, userGroup); err != nil {
		return nil, err
	}
	userGroups[userGroup.GetName()] = userGroup
	return userGroups, nil
}

func (b BackendDynamoDB) FindUserByName(name string) (map[string]model.UserGroup, error) {
	item, err := b.findByName(b.userTable, "user", name)
	if err != nil {
		return nil, err
	}

	return b.unmarshalUserGroup(new(model.User), map[string]model.UserGroup{}, item)
}

func (b BackendDynamoDB) FindUserByID(id int) (map[string]model.UserGroup, error) {
	item, err := b.findByID(b.userTable, "user", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	return b.unmarshalUserGroup(new(model.User), map[string]model.UserGroup{}, item)
}

func (b BackendDynamoDB) Users() (map[string]model.UserGroup, error) {
	items, err := b.findAll(b.userTable, "user")
	if err != nil {
		return nil, err
	}

	us := map[string]model.UserGroup{}
	for _, item := range items {
		if _, err := b.unmarshalUserGroup(new(model.User), us, item); err != nil {
			return nil, err
		}
	}
	return us, nil
}

func (b BackendDynamoDB) FindGroupByName(name string) (map[string]model.UserGroup, error) {
	item, err := b.findByName(b.groupTable, "group", name)
	if err != nil {
		return nil, err
	}

	return b.unmarshalUserGroup(new(model.Group), map[string]model.UserGroup{}, item)
}

func (b BackendDynamoDB) FindGroupByID(id int) (map[string]model.UserGroup, error) {
	item, err := b.findByID(b.groupTable, "group", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	return b.unmarshalUserGroup(new(model.Group), map[string]model.UserGroup{}, item)
}

func (b BackendDynamoDB) Groups() (map[string]model.UserGroup, error) {
	items, err := b.findAll(b.groupTable, "group")
	if err != nil {
		return nil, err
	}

	us := map[string]model.UserGroup{}
	for _, item := range items {
		if _, err := b.unmarshalUserGroup(new(model.Group), us, item); err != nil {
			return nil, err
		}
	}
	return us, nil

}

func (b BackendDynamoDB) highlowID(table string, high bool) int {
	ret := 0
	idAttr := "id"
	projExpr := "#ID"
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]string{
			"#ID": idAttr,
		},
		ProjectionExpression: &projExpr,
		TableName:            &table,
	}

	result, err := b.db.Scan(context.Background(), input)
	if err != nil {
		return ret
	}
	if len(result.Items) == 0 {
		return ret
	}
	for _, item := range result.Items {
		if idVal, ok := item["id"].(*types.AttributeValueMemberN); ok {
			id, err := strconv.Atoi(idVal.Value)
			if err != nil {
				return 0
			}
			if ret == 0 || (high && id > ret) || (!high && id < ret) {
				ret = id
			}
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
	av, err := attributevalue.MarshalMap(v)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &table,
	}

	_, err = b.db.PutItem(context.Background(), input)
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
	idStr := strconv.Itoa(id)
	params := &dynamodb.DeleteItemInput{
		TableName: &table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{
				Value: idStr,
			},
		},
	}
	_, err := b.db.DeleteItem(context.Background(), params)
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
