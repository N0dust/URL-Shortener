# URL-Shortener api

## Example
### create
POST http://127.0.0.1:8080/api/shorten
json {"url" : "https://https://www.xxxx.com","expiration_in_mainutes" : 2}
respond
json {"shortlink": "F"}

### get shorten information
GET http://127.0.0.1:8080/api/info?shortlink=F

### redirect
browser http://127.0.0.1:8080/F

