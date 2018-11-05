var params = {
    "TableName": "BTC_30m",
    "KeySchema": [
      { "AttributeName": "closeTime", "KeyType": "HASH" }
    ],
    "AttributeDefinitions": [
      { "AttributeName": "closeTime", "AttributeType": "N" }
    ],
    "ProvisionedThroughput": {
      "ReadCapacityUnits": 5,
      "WriteCapacityUnits": 5
    }
};
dynamodb.createTable(params, function(err, data) {
    if (err) ppJson(err); // an error occurred
    else ppJson(data); // successful response

});
dynamodb.listTables(params, function(err, data) {
    if (err) ppJson(err); // an error occurred
    else ppJson(data); // successful response
});