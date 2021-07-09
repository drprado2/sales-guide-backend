test:
	- go test -race ./...

test-cover:
	- go test -race -coverprofile cover.out ./... && go tool cover -html=cover.out -o cover.html && xdg-open ./cover.html

clear-cover:
	- sudo rm -rf cover.html cover.out

localstack-up:
	- cd eng/localstack && TMPDIR=/private$TMPDIR docker-compose up

localstack-up-force:
	- cd eng/localstack && docker-compose up --force-recreate -V

create-localstack-terraform-bucket-state:
	- aws --endpoint-url http://localhost:4566 s3api create-bucket --bucket terraform-state --region sa-east-1

init-terraform-local:
	- cd terraform/local && terraform init

local-tf-plan:
	- cd terraform/local && terraform plan -lock=false

local-tf-apply:
	- cd terraform/local && terraform apply -lock=false -auto-approve

aws-list-s3:
	- aws --endpoint-url http://localhost:4566 s3 ls --region sa-east-1

aws-list-dynamo-tables:
	- aws --endpoint-url http://localhost:4566 dynamodb list-tables --region sa-east-1

aws-list-api-gateway:
	- aws --endpoint-url http://localhost:4566 apigateway get-rest-apis --region sa-east-1

delete-localstack:
	- cd eng/localstack && docker-compose down

sns-add-email-sub:
	- aws sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:sns-test --protocol email --notification-endpoint mariaaug222@gmail.com --endpoint-url http://localhost:4566

sns-add-sqs-sub:
	- aws sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:sns-test --protocol sqs --notification-endpoint arn:aws:sqs:sa-east-1:000000000000:test_queue --endpoint-url http://localhost:4566

create-iam-user:
	- aws iam create-user --user-name sales-guide --endpoint-url http://localhost:4566 --region sa-east-1

publish-sns-test:
	- aws sns publish --topic-arn arn:aws:sns:sa-east-1:000000000000:sns-test --message test --endpoint-url http://localhost:4566

add-test-secretmanager:
	- cd eng/secretmanager && aws secretsmanager put-secret-value --secret-id test-sm --secret-string file://test-secret-values.json --endpoint-url http://localhost:4566 --region sa-east-1

get-test-secretmanager:
	- cd eng/secretmanager && aws secretsmanager get-secret-value --secret-id test-sm --version-stage AWSCURRENT --endpoint-url http://localhost:4566 --region sa-east-1

list-sfn:
	- aws stepfunctions list-state-machines --endpoint-url http://localhost:8083 --region sa-east-1

upload-hello-world-s3:
	- cd lambda/helloworld &&\
		GOOS=linux go build hello_world.go &&\
		zip hello_wolrd.zip hello_world &&\
		aws s3 cp hello_wolrd.zip s3://sandpit-sample/v1.0.2/hello_wolrd.zip --endpoint-url http://localhost:4566 --region sa-east-1 &&\
		rm hello_world &&\
		rm hello_wolrd.zip

upload-good_by-s3:
	- cd lambda/goodby &&\
		GOOS=linux go build good_by.go &&\
		zip good_by.zip good_by &&\
		aws s3 cp good_by.zip s3://sandpit-sample/v1.0.0/good_by.zip --endpoint-url http://localhost:4566 --region sa-east-1 &&\
		rm good_by &&\
		rm good_by.zip

upload-api_gateway-s3:
	- cd lambda/apiGatewayTest &&\
		GOOS=linux go build api_gateway.go &&\
		zip api_gateway.zip api_gateway &&\
		aws s3 cp api_gateway.zip s3://sandpit-sample/v1.0.1/api_gateway.zip --endpoint-url http://localhost:4566 --region sa-east-1 &&\
		rm api_gateway &&\
		rm api_gateway.zip

call-hello-world-lambda:
	- aws lambda invoke --function-name HelloWorld --endpoint-url http://localhost:4566 --region sa-east-1 --payload '{ "name": "Adriano", "age": 27 }' --cli-binary-format raw-in-base64-out hello_world_resp.json &&\
  	  cat hello_world_resp.json &&\
  	  rm hello_world_resp.json

call-good_by-lambda:
	- aws lambda invoke --function-name GoodBy --endpoint-url http://localhost:4566 --region sa-east-1 --payload '{ "name": "Adriano", "age": 27 }' --cli-binary-format raw-in-base64-out good_by_resp.json &&\
  	  cat good_by_resp.json &&\
  	  rm good_by_resp.json

list-redis-groups:
	- aws elasticache describe-replication-groups \
      --endpoint-url http://localhost:4566 \
      --region sa-east-1