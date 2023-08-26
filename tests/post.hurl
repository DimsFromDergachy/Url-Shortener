POST http://localhost:8082/
{
	"url": "https://duckduckgo.com",
	"alias": "test_alias"
}

HTTP 200
[Asserts]
jsonpath "$.status" == "OK"
jsonpath "$.error" not exists
jsonpath "$.alias" == "test_alias"