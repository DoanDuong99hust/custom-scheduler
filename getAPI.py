import requests
import json
from types import SimpleNamespace

url = "http://localhost:5000/api/v1/test"
response = requests.get(url)
data = response.text
parsed = json.loads(data)
# raw_data = json.loads(parsed, object_hook=lambda d: SimpleNamespace(**d))

print(parsed['switch'])
for i in parsed['switch']:
    for j in i['port']:
        print(j['Rx'])