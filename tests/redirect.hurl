GET http://localhost:8082/test_alias

HTTP 302
[Asserts]
header "Location" == "https://duckduckgo.com"