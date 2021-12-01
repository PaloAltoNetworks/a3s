#!/bin/python

import json
import ssl
from flask import Flask

app = Flask(__name__)


@app.route("/mod")
def modifyclaims():
    return json.dumps(["hello=world"])


if __name__ == "__main__":
    ssl_context = ssl.create_default_context(
        purpose=ssl.Purpose.CLIENT_AUTH,
        cafile='certs/ca-cert.pem'
    )

    ssl_context.load_cert_chain(
        certfile='certs/claimmod-cert.pem',
        keyfile='certs/claimmod-key.pem',
    )

    ssl_context.verify_mode = ssl.CERT_REQUIRED

app.run(
    port=5001,
    ssl_context=ssl_context,
)
