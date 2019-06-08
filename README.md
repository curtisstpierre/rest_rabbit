# rest_rabbit
A rest api rabbit producer

# Getting Started
```
go get
go build
./rest_rabbit
```

# Testing in Python

```
import requests

params = {
    'grant_type': 'client_credentials',
    'client_id':'foo',
    'client_secret':'bar',
    'scope':'read'
}

x=requests.post('http://localhost:8000/token', data=params).json()
token = x['access_token']
headers = { 'Authorization': f'Bearer {token}'}
requests.post(f'http://localhost:8000/messages/foobar', headers=headers)
```
