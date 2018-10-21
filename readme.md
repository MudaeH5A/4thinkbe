# 4think backend

## Routes

- GET /:userNumber
    - registers new user with predefined params
- POST /:userNumber/:vehicle
    - after boxes analisys, a post to this URI updates the payment information
- GET /:userNumber/:room/:boxNumber/code
    - generates QR codes for a box
- GET /:userNumber/:room/:boxNumber
    - URI that a QR code shows to user after being scanned.


## Tech Stack

- Language: Go go1.10.2 linux/amd64
- Frameworks and libraries directly used: 
    - Echo (routing, github.com/labstack/echo)
    - Go-QRCode (github.com/skip2/go-qrcode)
    - Mongo driver (gopkg.in/mgo.v2)
- Deploy:
    - Heroku
- Backend: https://mudae.herokuapp.com/7
- Frontend: http://www.mudae.com.br/