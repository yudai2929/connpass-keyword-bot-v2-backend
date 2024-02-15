package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
	"github.com/yudai2929/connpass-keyword-bot-v2-backend/cdk/pkg/lambda"
	"github.com/yudai2929/connpass-keyword-bot-v2-backend/cdk/pkg/utils"
)

type ConnpassBotV2Props struct {
	awscdk.StackProps
}

func NewConnpassBotV2(scope constructs.Construct, id string, props *ConnpassBotV2Props) (awscdk.Stack, error) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	apiLambdaConfig := lambda.NewConfig(
		"bin/api",
		"app/cmd/api/main.go",
		"ConnpassBotV2ApiFunction",
	)

	if err := utils.GolangBuild(apiLambdaConfig.BuildOutputPath, apiLambdaConfig.GolangPath); err != nil {
		return stack, err
	}

	lambdaFn := awslambda.NewFunction(stack, jsii.String(apiLambdaConfig.FunctionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(apiLambdaConfig.FunctionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(apiLambdaConfig.BuildDir), nil),
		Handler:      jsii.String(apiLambdaConfig.Handler),
	})

	api := awsapigateway.NewRestApi(stack, jsii.String("ConnpassBotV2ApiGateway"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("ConnpassBotV2ApiGateway"),
	})

	apiResource := api.Root().AddResource(jsii.String("api"), nil)

	proxyResource := apiResource.AddResource(jsii.String("{proxy+}"), nil)
	proxyResource.AddMethod(jsii.String("ANY"), awsapigateway.NewLambdaIntegration(lambdaFn, nil), nil)

	return stack, nil
}

func main() {
	defer jsii.Close()

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	app := awscdk.NewApp(nil)

	_, err := NewConnpassBotV2(app, "ConnpassBotV2", &ConnpassBotV2Props{
		awscdk.StackProps{
			Env: env(),
		},
	})

	if err != nil {
		panic(err)
	}

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT")),
		Region:  jsii.String(os.Getenv("AWS_REGION")),
	}
}
