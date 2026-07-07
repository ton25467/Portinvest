import urllib.request
import urllib.error
import json

def get_token():
    req = urllib.request.Request("http://localhost:8080/api/v1/auth/login", data=json.dumps({"email":"user@example.com", "password":"password"}).encode(), headers={"Content-Type":"application/json"})
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read())['token']

token = get_token()

def get_api(path):
    req = urllib.request.Request("http://localhost:8080/api/v1" + path, headers={"Authorization": f"Bearer {token}"})
    try:
        with urllib.request.urlopen(req) as resp:
            print(f"--- {path} ---")
            print(json.dumps(json.loads(resp.read()), indent=2))
    except Exception as e:
        print(f"Error {path}:", e)

get_api("/portfolios")
get_api("/portfolios/00000000-0000-0000-0000-000000000002/summary")
get_api("/portfolios/00000000-0000-0000-0000-000000000002/holdings")
