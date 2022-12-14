package main

import (
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	lambdanodejs "github.com/aws/aws-cdk-go/awscdk/v2/awslambdanodejs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ExternalStackProps struct {
	cdk.StackProps
}

func NewExternalStack(scope constructs.Construct, id string, props *ExternalStackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	externalApi := apigateway.NewRestApi(stack, jsii.String("externalApi"), &apigateway.RestApiProps{})

	getProductFn := lambdanodejs.NewNodejsFunction(
		stack,
		jsii.String("getProductFn"),
		&lambdanodejs.NodejsFunctionProps{
			Handler: jsii.String("handler"),
			Entry:   jsii.String("./src/lambdas/get-product.ts"),
			Runtime: lambda.Runtime_NODEJS_16_X(),
			Bundling: &lambdanodejs.BundlingOptions{
				Minify: jsii.Bool(true),
			},
		},
	)

	api := externalApi.Root().AddResource(jsii.String("api"), &apigateway.ResourceOptions{})
	productApi := api.AddResource(jsii.String("{id}"), &apigateway.ResourceOptions{})
	productApi.AddMethod(
		jsii.String("GET"),
		apigateway.NewLambdaIntegration(getProductFn, &apigateway.LambdaIntegrationOptions{}),
		api.DefaultMethodOptions(),
	)

	cdk.NewCfnOutput(stack, jsii.String("externalApiUrl"), &cdk.CfnOutputProps{
		Value:      externalApi.Url(),
		ExportName: jsii.String("externalApiUrl"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := cdk.NewApp(nil)

	NewExternalStack(app, "ExternalStack", &ExternalStackProps{
		cdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *cdk.Environment {
	return &cdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
