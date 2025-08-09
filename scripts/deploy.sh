#!/bin/bash

# Script para implantar a aplicação Go no LocalStack

set -e

# Variáveis de configuração
LOCALSTACK_ENDPOINT="http://localhost:4566"
AWS_REGION="us-east-1"
FUNCTION_NAME="myLambdaFunction"

# Criar a função Lambda
echo "Criando a função Lambda..."
aws --endpoint-url $LOCALSTACK_ENDPOINT lambda create-function \
  --function-name $FUNCTION_NAME \
  --runtime go1.x \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --handler main \
  --zip-file fileb://function.zip

# Criar o API Gateway
echo "Criando o API Gateway..."
API_ID=$(aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway create-rest-api --name "MyAPI" --query 'id' --output text)

# Obter o ID do recurso raiz
ROOT_RESOURCE_ID=$(aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway get-resources --rest-api-id $API_ID --query 'items[?path==`/`].id' --output text)

# Criar um novo recurso
RESOURCE_ID=$(aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_RESOURCE_ID --path-part "myresource" --query 'id' --output text)

# Criar um método para o recurso
aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method POST --authorization-type "NONE"

# Integrar o método com a função Lambda
aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway put-integration --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method POST --type AWS_PROXY --integration-http-method POST --uri "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:$FUNCTION_NAME/invocations"

# Implantar a API
aws --endpoint-url $LOCALSTACK_ENDPOINT apigateway create-deployment --rest-api-id $API_ID --stage-name dev

echo "Implantação concluída!"