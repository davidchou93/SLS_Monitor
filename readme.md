# A group of serverless TG-bot services 
For private usage only, functions are not guaranteed.

*Due to AWS serverless api change, some functions are no longer working*

## Functions
### Text to audio
Using AWS-Polly to converse input text to audio file.

PATH: `/deliever`

### Crypto price monitor
Monitor selected crypto price every 30 minutes, notice if price changed more than 2%.

PATH: `/echo`

### DynamoDB receiver
Put the given input(request body) into DynamoDB.

PATH: `/receiver`

## Toolbox
In order to make local development environment, following tools are needed.
- [AWS Serverless Application Model (SAM)](https://github.com/awslabs/serverless-application-model)
- [Serverless framework(serverless)](https://serverless.com/)
- [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)
- [Golang package management tool(Dep)](https://github.com/golang/dep)(Not recommended for future development, use `mod` instead)