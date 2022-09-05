package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/drprado2/sales-guide/pkg/awsconfig"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/jinzhu/copier"
	"reflect"
)

var (
	client *dynamodb.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = dynamodb.NewFromConfig(cfg)
	return nil
}

func ListTablesSvc(ctx context.Context) ([]string, error) {
	resp, err := client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		logs.Logger(ctx).Infof("failed to list dynamo tables, %v", err)
		return nil, err
	}

	return resp.TableNames, nil
}

func GetItemSvc(
	ctx context.Context,
	table string,
	keys map[string]interface{},
	projectionExpression string,
	target interface{}) error {

	mKeys := make(map[string]types.AttributeValue)
	for k, v := range keys {
		value, err := attributevalue.Marshal(v)
		if err != nil {
			logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
			return err
		}
		mKeys[k] = value
	}
	ipt := &dynamodb.GetItemInput{
		Key:                  mKeys,
		TableName:            &table,
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: &projectionExpression,
	}
	resp, err := client.GetItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed get item in dynamo tables, %v", err)
		return err
	}
	if resp.Item == nil {
		return ItemNotFoundError
	}

	if err := attributevalue.UnmarshalMap(resp.Item, target); err != nil {
		logs.Logger(ctx).Errorf("failed unmarshal result err=%v", err)
		return err
	}
	return nil
}

func PutItemSvc(
	ctx context.Context,
	table string,
	item interface{}) error {

	mItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
		return err
	}
	ipt := &dynamodb.PutItemInput{
		Item:      mItem,
		TableName: &table,
	}
	_, err = client.PutItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed to put item in dynamo, %v", err)
		return err
	}
	return nil
}

func DeleteItemSvc(
	ctx context.Context,
	table string,
	keys map[string]interface{}) error {

	mKeys := make(map[string]types.AttributeValue)
	for k, v := range keys {
		value, err := attributevalue.Marshal(v)
		if err != nil {
			logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
			return err
		}
		mKeys[k] = value
	}
	ipt := &dynamodb.DeleteItemInput{
		Key:       mKeys,
		TableName: &table,
	}
	_, err := client.DeleteItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed to list dynamo tables, %v", err)
		return err
	}
	return nil
}

func UpdateItemSvc(
	ctx context.Context,
	table string,
	keys map[string]interface{},
	updates map[string]interface{}) error {

	mKeys := make(map[string]types.AttributeValue)
	for k, v := range keys {
		value, err := attributevalue.Marshal(v)
		if err != nil {
			logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
			return err
		}
		mKeys[k] = value
	}
	mUpdates := make(map[string]types.AttributeValueUpdate)
	for k, v := range updates {
		value, err := attributevalue.Marshal(v)
		if err != nil {
			logs.Logger(ctx).Errorf("error in marshal attribute at dynamo UpdateItem err=%v", err)
			return err
		}
		mUpdates[k] = types.AttributeValueUpdate{
			Action: "PUT",
			Value:  value,
		}
	}

	ipt := &dynamodb.UpdateItemInput{
		Key:              mKeys,
		TableName:        &table,
		AttributeUpdates: mUpdates,
	}
	_, err := client.UpdateItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed to update item in dynamo, %v", err)
		return err
	}
	return nil
}

func BatchGetItemSvc(
	ctx context.Context,
	table string,
	filters []map[string]interface{},
	projectionExpression string,
	target interface{}) ([]interface{}, error) {

	mFilters := make([]map[string]types.AttributeValue, len(filters))
	for i, filter := range filters {
		mKeys := make(map[string]types.AttributeValue)
		for k, v := range filter {
			value, err := attributevalue.Marshal(v)
			if err != nil {
				logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
				return nil, err
			}
			mKeys[k] = value
		}
		mFilters[i] = mKeys
	}
	keysAndValues := types.KeysAndAttributes{
		Keys:                 mFilters,
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: &projectionExpression,
	}
	reqItems := make(map[string]types.KeysAndAttributes)
	reqItems[table] = keysAndValues

	ipt := &dynamodb.BatchGetItemInput{
		RequestItems:           reqItems,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}
	resp, err := client.BatchGetItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed get item in dynamo tables, %v", err)
		return nil, err
	}
	if resp.Responses == nil || len(resp.Responses) == 0 {
		return nil, ItemNotFoundError
	}

	result := make([]interface{}, 0)

	for _, item := range resp.Responses[table] {
		if err := attributevalue.UnmarshalMap(item, target); err != nil {
			logs.Logger(ctx).Errorf("failed unmarshal result err=%v", err)
			return nil, err
		}
		current := reflect.Indirect(reflect.ValueOf(target)).Interface()
		copier.Copy(target, current)
		result = append(result, current)

	}
	return result, nil
}

func BatchPutItemsSvc(
	ctx context.Context,
	table string,
	elements []interface{}) error {

	writes := make([]types.WriteRequest, len(elements))
	for i, ele := range elements {
		val, err := attributevalue.MarshalMap(ele)
		if err != nil {
			logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
			return err
		}
		writes[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: val,
			},
		}
	}

	reqItems := make(map[string][]types.WriteRequest)
	reqItems[table] = writes
	ipt := &dynamodb.BatchWriteItemInput{
		RequestItems:                reqItems,
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
	}
	_, err := client.BatchWriteItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed batch write items in dynamo, %v", err)
		return err
	}
	return nil
}

func BatchDeleteItemsSvc(
	ctx context.Context,
	table string,
	keys []map[string]interface{}) error {

	filters := make([]map[string]types.AttributeValue, 0)
	writes := make([]types.WriteRequest, len(keys))
	for i, kk := range keys {
		mKeys := make(map[string]types.AttributeValue)
		for k, v := range kk {
			value, err := attributevalue.Marshal(v)
			if err != nil {
				logs.Logger(ctx).Errorf("error in marshal attribute at dynamo GetItem err=%v", err)
				return err
			}
			mKeys[k] = value
		}

		writes[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: mKeys,
			},
		}
		filters = append(filters, mKeys)
	}

	reqItems := make(map[string][]types.WriteRequest)
	reqItems[table] = writes
	ipt := &dynamodb.BatchWriteItemInput{
		RequestItems:                reqItems,
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
	}
	_, err := client.BatchWriteItem(ctx, ipt)
	if err != nil {
		logs.Logger(ctx).Errorf("failed batch write items in dynamo, %v", err)
		return err
	}
	return nil
}
