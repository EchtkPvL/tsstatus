# Luzifer / tsstatus

`tsstatus` is a small utility to expose a status of a TeamSpeak3 server.

This can be used to have a monitoring for the server (HTTP 200 vs. HTTP 500) and to retrieve a list of users being present in a channel

```console
# curl -sSf http://localhost:3000/status | jq .
{
  "info": {
    "server": {
      "clients_online": 2,
      "query_clients_online": 1,
      "max_clients": 512,
      "total_channels": 37,
      "name": "EchtkPvL.de",
      "port": 9987,
      "status": "online",
      "uptime": 18704453,
      "version": "3.13.6 [Build: 1623234157]",
      "platform": "Linux"
    },
    "channels": [{
        "id": 1074,
        "name": "Eingangshalle",
        "clients": [{
            "away": false,
            "away_message": "",
            "nickname": "EchtkPvL"
          }
        ]
      }
    ]
  }
}
```

In case of failure

```console
# curl -sSf http://localhost:3000/status | jq .
{
  "error": "Unable to create client: dial tcp 127.0.0.1:10011: connect: connection refused"
}
```

## Run

### Clone
`git clone git@github.com:EchtkPvL/tsstatus.git && cd tsstatus/`

### Build
`docker build -t tsstatus .`

### Run once
`docker run -p 3000:3000 tsstatus --query-user queryuser --query-pass PASSW0RD --server-address tsserver.tld --server-id 1`

### Run daemonized
`docker run -d -p 3000:3000 --restart unless-stopped --name tsstatus tsstatus --query-user queryuser --query-pass PASSW0RD --server-address tsserver.tld --server-id 1`
