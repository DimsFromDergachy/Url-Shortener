# There is no check of existing
DELETE http://localhost:8082/test_auth

HTTP 200

POST http://localhost:8082/
{
	"url": "https://duckduckgo.com",
	"alias": "test_auth"
}

HTTP 403