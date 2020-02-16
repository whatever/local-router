# Routing for fun

Worse than HA Proxy. Route traffic based on request URI's Host to locally running webservers.


## Run

`sudo go run router.go -port 8080 -config whatever.config`


## Config file

New-line separated in the format `%s: %d`

```
whatever.org: 9000
nevermind.business: 9001
anyway.net:  9002
```

