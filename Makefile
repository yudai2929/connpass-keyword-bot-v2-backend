run-api:
	@echo "Running API"
	@go run app/cmd/api/main.go


cdk-deploy:
	@echo "Deploying CDK"
	@cdk deploy

create-sam-template:
	@echo "Creating SAM template"
	@cdk synth --no-staging > template.yaml

PHONY: run-api