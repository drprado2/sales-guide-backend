test:
	- go test -race ./...

apply-flyway:
	- flyway -configFiles=./flyway-structure.conf migrate &&\
 	 	flyway -configFiles=./flyway-data.conf migrate

export-aws-env-vars:
	- export AWS_ACCESS_KEY_ID=AKIA4TOL2AU6WVIAJABA &&\
      export AWS_SECRET_ACCESS_KEY=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 &&\
      export AWS_DEFAULT_REGION=us-east-2

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

init-terraform-aws:
	- cd terraform/mariaaug222 && terraform init

aws-tf-plan:
	- cd terraform/mariaaug222 && terraform plan -lock=false

aws-tf-apply:
	- cd terraform/mariaaug222 && terraform apply -lock=false -auto-approve

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

build-custom-kong-img:
	- cd eng/kong && docker image build -t custom-kong-ee .

create-kong-db:
	- docker pull kong/kong-gateway:2.5.0.0-alpine &&\
		docker tag kong/kong-gateway:2.5.0.0-alpine kong-ee &&\
		docker network create kong-ee-net &&\
		docker run -d --name kong-ee-database \
          --network=kong-ee-net \
          -p 25432:5432 \
          -e "POSTGRES_USER=kong" \
          -e "POSTGRES_DB=kong" \
          -e "POSTGRES_PASSWORD=kong" \
          postgres:9.6

migrate-kong-db:
	-  docker run --rm --network=kong-ee-net \
          -e "KONG_DATABASE=postgres" \
          -e "KONG_PG_HOST=kong-ee-database" \
          -e "KONG_PG_PORT=5432" \
          -e "KONG_PG_PASSWORD=kong" \
          -e "KONG_PASSWORD=Admin123!" \
          kong-ee kong migrations bootstrap

start-kong:
	- docker run -d --name kong-ee --network=kong-ee-net \
	  -e "KONG_DATABASE=postgres" \
	  -e "KONG_PG_HOST=kong-ee-database" \
	  -e "KONG_PG_PASSWORD=kong" \
	  -e "KONG_PG_PORT=5432" \
	  -e "KONG_PLUGINS=bundled,kong-jwt2header" \
	  -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
	  -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
	  -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
	  -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
	  -e "KONG_ADMIN_LISTEN=0.0.0.0:8001" \
	  -e "KONG_ADMIN_GUI_URL=http://localhost:8002" \
		-p 8000:8000 \
		-p 8443:8443 \
		-p 8001:8001 \
		-p 8444:8444 \
		-p 8002:8002 \
		-p 8445:8445 \
		-p 8003:8003 \
		-p 8004:8004 \
		custom-kong-ee

clear-kong:
	- docker container rm -f kong-ee &&\
 		docker container rm -f kong-ee-database &&\
 		docker container rm -f kong &&\
 		docker network rm kong-ee-net

sync-kong:
	- cd eng/kong && deck sync

run-local-dependencies:
	- cd eng && docker-compose up -d
