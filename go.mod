module github.com/drprado2/react-redux-typescript

go 1.15

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go-v2 v1.7.0
	github.com/aws/aws-sdk-go-v2/config v1.4.1
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.1.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.4.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.4.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.4.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.6.0
	github.com/aws/smithy-go v1.5.0
	github.com/felixge/fgprof v0.9.1
	github.com/go-redis/redis/v8 v8.11.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/gosidekick/goconfig v1.3.0
	github.com/jackc/pgx/v4 v4.10.1
	github.com/jinzhu/copier v0.3.2
	github.com/lib/pq v1.9.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.0
)
