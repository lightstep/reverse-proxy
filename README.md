# reverse-proxy

reverse-proxy is a simple tool that takes an incoming request and sends it to another server, proxying the response back to the client. It is used for connecting with the HTTPS endpoint of the LightStep public satellite. 

## Instructions

To start the proxy with defaults, run:
```
docker run -p 8126:8126 lightstep/reverse-proxy:latest
```

To view the complete list of options run the image with the `help` flag.
```
docker run lightstep/reverse-proxy:latest --help
```

## Flags
```
 -ca string
        Custom CA certificate
 -forward-url string
        Satellite address (default "https://ingest.lightstep.com/")
 -port string
        Port for the proxy server (default "8126")
```

