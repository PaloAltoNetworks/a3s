#!/bin/python

import json
import ssl
from flask import Flask, request, Response

app = Flask(__name__)


@app.route("/login", methods=["POST"])
def login():
    creds = request.get_json()
    if creds["username"] != "john" or creds["password"] != "pass":
        return Response("Forbidden", 401, {})
    return json.dumps(["user=john"])


if __name__ == "__main__":
    ssl_context = ssl.create_default_context(
        purpose=ssl.Purpose.CLIENT_AUTH,
        cafile='certs/ca-cert.pem'
    )

    ssl_context.load_cert_chain(
        certfile='certs/httpsource-cert.pem',
        keyfile='certs/httpsource-key.pem',
    )

    ssl_context.verify_mode = ssl.CERT_REQUIRED

app.run(
    port=5002,
    ssl_context=ssl_context,
)
