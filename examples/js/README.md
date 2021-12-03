# JS Demo

A simple UI that handles authentication using a3s. Currently it supports `MTLS`, `LDAP` and `OIDC`. 

## Install requirements

You need to have [nodejs](https://nodejs.org/) and [yarn (v1)](https://yarnpkg.com/getting-started/install) installed.

## Launch the app

- For the first time ony, run `yarn` to install the dependencies.
- Then, run `yarn start` to start the local dev server.
- Follow the instruction to open the app in your browser.
> If `yarn start` fails, stop it and run again. There seems to be some bug with snowpack. We'll resolve it soon.
> For CORS issue, start Chrome using `open -n -a /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --args --user-data-dir="/tmp/chrome_dev_test" --disable-web-security` (this command works for macOS only).
> For cert issue, open the a3s server url in your broswer and trust the certificate.