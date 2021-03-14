# URL-Shortener api

## Example


### Create
POST http://127.0.0.1:8080/api/shorten


json {"url" : "https://https://www.xxxx.com","expiration_in_mainutes" : 2}


get respond  {"shortlink": "F"}

### Get shorten information
GET http://127.0.0.1:8080/api/info?shortlink=F

### Redirect
browser http://127.0.0.1:8080/F

