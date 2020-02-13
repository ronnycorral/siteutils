import logging
import boto3
import pymysql
import os
import sys

logger = logging.getLogger()
logger.setLevel(logging.INFO)

s3 = boto3.client('s3')

host = os.environ['host']
user = os.environ['user']
passwd = os.environ['passwd']
db = os.environ['db']

def lambda_handler(event, context):

    try:
        conn = pymysql.connect(host, user=user, passwd=passwd, db=db, connect_timeout=5)
        cur = conn.cursor()
    except pymysql.MySQLError as e:
        logger.error("ERROR: Unexpected error: Could not connect to MySQL instance.")
        logger.error(e)
        sys.exit()

    # retrieve bucket name and file_key from the S3 event
    bucket_name = event['Records'][0]['s3']['bucket']['name']
    file_key = event['Records'][0]['s3']['object']['key']
    logger.info('Reading {} from {}'.format(file_key, bucket_name))
    # get the object
    obj = s3.get_object(Bucket=bucket_name, Key=file_key)
    lines = obj['Body'].read().strip().split(b'\n')
    for line in lines:
        #logger.info(line.decode())
        website_values = line.decode().split('\t')
        mysql = "INSERT INTO logdata (ip,requestDate,request,responseCode,CountryCode,regionCode) VALUES (%s, STR_TO_DATE(%s,'%%d/%%b/%%Y:%%T'), %s, %s, %s, %s)"
        cur.execute(mysql,website_values)
    conn.commit()
