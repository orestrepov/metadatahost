# Metadata Host Service
Metadata Host Web Service in GO

It is proof of concept of developing a web application in "Golang" to search for information associated with domains and the history of recent domains searches.

# API Service
# Build 
go build -o metadatahost .

# Run serve
./metadatahost serve

# Endpoints
##### Example to search domain data:
curl -v http://localhost:9095/api/hosts/search/www.mydomain.com
##### Example to get history of domain search:
curl -v http://localhost:9095/api/hosts/search/history
