import logging
import boto3
import re
import requests
import json
from xml.dom.minidom import parseString

#GEOIP related
geo_url = "https://api.ipdata.co/%s?api-key=%s&fields=country_code,region_code"
country_codes = {}
region_codes = {}

regex = '([(\d\.)]+) - - \[(\S*) \+0000\] "GET (.*?) HTTP/1.\d" (\d+) ([0-9,-]+) "(.*?)" "(.*?)"'
excludes = (
    b'.gif ',          # don't care about counting images, css or javascript files
    b'.jpg ',
    b'.png ',
    b'.ico ',
    b'.css ',
    b'.js ',
    b'robots.txt ',   # don't care since they are robots
    b'54.172.87.69',  # corral.com ip
    b'76.126.34.55X'  # me at home
    )

logger = logging.getLogger()
logger.setLevel(logging.INFO)

s3 = boto3.client('s3')
ssm = boto3.client('ssm')

def lambda_handler(event, context):

    geoip_key = ssm.get_parameter(Name='geoipkey', WithDecryption=False)['Parameter']['Value']
    log_data_output = ""
    # retrieve bucket name and file_key from the S3 event
    bucket_name = event['Records'][0]['s3']['bucket']['name']
    file_key = event['Records'][0]['s3']['object']['key']
    new_file_key = file_key.replace('WebsiteLogs','ParsedLogData')
    logger.info('Reading {} from {}'.format(file_key, bucket_name))
    # get the object
    obj = s3.get_object(Bucket=bucket_name, Key=file_key)
    lines = obj['Body'].read().split(b'\n')
    for line in lines:
        # logger.info(line.decode())
        if not any(s in line for s in excludes):
            result =  re.match(regex, line.decode('utf-8'))
            if result:
                log_line_info =  result.groups()
                ip = log_line_info[0]
                request_date = log_line_info[1]
                request = log_line_info[2].replace("'","\\'")[:250]
                response_code = log_line_info[3]
                if ip not in country_codes:
                    r = requests.get(geo_url % (ip, geoip_key))
                    json_content = json.loads(r.content)
                    try:
                        country_codes[ip] = json_content["country_code"]
                    except:
                        country_codes[ip] = 'Unknown'
                    try:
                        region_codes[ip] = json_content["region_code"]
                    except:
                        region_codes[ip] = 'Unknown'
                log_data_output += "%s\t%s\t%s\t%s\t%s\t%s\n" % (ip,request_date,request,response_code,country_codes[ip],region_codes[ip])
                #logger.info(logDataOutput)

    if len(log_data_output):
        s3.put_object(Bucket=bucket_name, Key=new_file_key, Body=log_data_output.encode("utf-8"))
