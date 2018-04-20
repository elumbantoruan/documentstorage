# Document Storage
simple document storage

## Dependencies:
* github.com/gorilla/mux
    * HTTP router
* github.com/dgrijalva/jwt-go
    * JWT Token 

## Project Structures: 
### credentials
It's a folder where user credentials are stored

### handler
It's a package to manage the handler for files, login, and register resources.

### model
It's a package for request and response payload

### repository
It's a package for interfaces: user credential and user storage. It contains the concrete implementation using file system and mock storage (in memory)

### storage
It's a folder where user storage are stored
