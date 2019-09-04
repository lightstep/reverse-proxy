# reverse-proxy

reverse-proxy is a simple tool that takes an incoming request and sends it to another server, proxying the response back to the client. It is used for connecting with the HTTPS endpoint of the LightStep public satellite. 

## Flags
```
 -ca string
        Custom CA certificate
 -forward-url string
        Satellite address (default "https://ingest.lightstep.com/")
 -port string
        Port for the proxy server (default "8126")
```

