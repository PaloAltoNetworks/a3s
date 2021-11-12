#!/bin/python

from flask import Flask, request, Response
from functools import wraps
import requests

app = Flask(__name__)

def authenticate(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        auth = request.authorization
        if not auth or not auth.username == "Bearer" or auth.password == "":
            return Response('Forbidden', 401, {})
        if requests.post(
            "https://127.0.0.1:44443/authz",
            verify=False,
            headers={'Content-Type': 'application/json'},
            json = {
                'token': auth.password,
                'action': request.method,
                'resource': request.path,
                'namespace': "/testapp",
                'audience': "testapp",
            },
        ).status_code != 204:
            return Response('Forbidden', 403, {})
        return f(*args, **kwargs)
    return wrapper

@app.route("/")
def public(): 
    return "This is public. try to access /secret or /topsecret"


@app.route("/secret")
@authenticate
def secret(): 
    return "This is secret! Noice!"

@app.route("/topsecret")
@authenticate
def topsecret(): 
    return "This is top secret! Awesome!"

if __name__ == "__main__":
    app.run(ssl_context='adhoc')
