test:
	- go test -race ./...

test-cover:
	- go test -race -coverprofile cover.out ./... && go tool cover -html=cover.out -o cover.html && xdg-open ./cover.html

clear-cover:
	- sudo rm -rf cover.html cover.out

localstack-up:
	- cd eng/localstack && TMPDIR=/private$TMPDIR docker-compose up

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

delete-localstack:
	- cd eng/localstack && docker-compose down
