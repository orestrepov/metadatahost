# metadatahost
Metadata Host Web Application in GO

It is proof of concept of developing a web application in Go to search for information associated with domains and the history of recent searches.

#API Service
#Build 
go build -o metadatahost .

#Run serve
./metadatahost serve

#Endpoints
#####Example to search domain data:
curl -v http://localhost:9095/api/hosts/search/www.mydomain.com
#####Example to get history of domain search:
curl -v http://localhost:9095/api/hosts/search/history
