Echo "[I] DynamoDb Local Starting... "
Start-Process Powershell -ArgumentList "java -D'java.library.path=./dynamoDb/DynamoDBLocal_lib' -jar ./dynamoDb/DynamoDBLocal.jar -sharedDb" 
Echo "[I] DynamoDb Local Started!"
aws dynamodb create-table --cli-input-json file://create_table.json --endpoint-url http://localhost:8000
Echo "[I] Show existed Tables"
aws dynamodb list-tables --endpoint-url http://localhost:8000
Echo "[I] SAM Starting..."
sam local start-api