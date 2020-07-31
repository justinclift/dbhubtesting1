Just a temporary repo for developing test code in for the new DBHub.io Go API

This code will likely be merged into the DBHub.io Go API, as part of its testing pieces.

To use this, you'll need to create a simple text file under `~/.dbhub/apiclient.toml`.

It should contain your DBHub.io API key, and (optionally) an alternative server
for the API client to talk to.

Example for doing local API development:

```
[api]
api_key = "YOUR_API_KEY_HERE"
server = "https://jctesting1.dbhub.io:8443"
```

Example contacting the real (production) API server:

```
[api]
api_key = "YOUR_API_KEY_HERE"
server = ""
```
