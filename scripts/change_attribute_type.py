import os

import boto3
from boto3.dynamodb.conditions import Key

db_client = boto3.resource('dynamodb')
table_res = db_client.Table(os.environ['TABLE_NAME'])



not_finished = True
ret = table_res.scan()
while not_finished:
    for item in ret['Items']:
        urls = item['urls']
        new_item = item
        if isinstance(urls, list):
            new_item['urls'] = set(urls)
            print(f'urls of {item["Title"]} is converted to set')
            table_res.put_item(Item = new_item)
    if "LastEvaluatedKey" in ret:
        last_key = ret['LastEvaluatedKey']
        ret = table_res.scan(ExclusiveStartKey = last_key)
    else:
        not_finished = False
