#!/bin/python

from flask import Flask, request, Response, redirect
from functools import wraps
import requests

app = Flask(__name__)
# CORS(app)

self_url = 'https://localhost:5000'
a3s_url = 'https://127.0.0.1:44443'


# authenticator:start
def authenticate(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        password = request.cookies.get("x-a3s-token")
        auth = request.authorization
        if auth and auth.username == "Bearer" and auth.password != "":
            password = auth.password
        if not password and request.args.get("rlogin") is not None:
            return redirect(
                '%s/login?proxy=%s&redirect=%s/%s&audience=%s' %
                (self_url, self_url, self_url, request.path, "testapp")
            )
        if requests.post(
            "%s/authz" % a3s_url,
            verify=False,
            headers={'Content-Type': 'application/json'},
            json={
                'token': password,
                'action': request.method,
                'resource': request.path,
                'namespace': "/testapp",
                'audience': "testapp",
            },
        ).status_code != 204:
            return Response('Forbidden\n', 403, {})
        return f(*args, **kwargs)
    return wrapper
# authenticator:end


# routes:start
@ app.route("/")
def public():
    return "This is public. try to access <a href=/secret?rlogin>/secret</a> or <a href=/topsecret?rlogin>/topsecret</a>\n"


@ app.route("/secret")
@ authenticate
def secret():
    return "This is secret! Noice!\n"


@ app.route("/topsecret")
@ authenticate
def topsecret():
    return "This is top secret! Awesome!\n"
# routes:end


@ app.route("/issue", methods=['POST'])
def issue():
    resp = requests.post(
        "%s/issue" % a3s_url,
        verify=False,
        json=request.get_json(silent=True),
        headers={'Content-Type': 'application/json'},
    )
    headers = [(name, value) for (name, value) in resp.raw.headers.items()]
    return Response(resp.content, resp.status_code, headers)


@ app.route("/login", methods=['GET'])
def login():
    return requests.get(
        "%s/ui/login.html" % a3s_url,
        verify=False,
        params=request.args,
        headers=request.headers,
    ).content


if __name__ == "__main__":
    app.run(ssl_context='adhoc', debug=True)
