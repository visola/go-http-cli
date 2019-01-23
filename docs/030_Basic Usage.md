## Basic Usage

You can use go-http-cli similarly as you use cURL. The most commonly used options are available and the same. So you can do something like:

```bash
$ http \
  -H Content-Type:application/json \
  -X POST \
  -d '{ "name": "John Doe" }' \
  https://httpbin.org/post?companyId=1234
```

Will execute the following:
```HTTP
POST https://httpbin.org/post?companyId=1234
Content-Type: application/json

{ "name": "John Doe" }
```

go-http-cli can help you to do URL encoding if you pass in some key-value pairs, instead of adding the query string to the end of the URL, like the following, which would generate the same as above:

```bash
http \
  -H Content-Type:application/json \
  -X POST \
  -d '{ "name": "John Doe" }' \
  https://httpbin.org/post \
  companyId=1234
```

This is equivalent of the above. Except it would automatically URL encode the values to put in the query string.

When dealing with REST APIs, it's common to build JSONs and sending them in the body. go-http-cli can do this for you. If key-value pairs are passed in and the method is `POST`, it will automatically build a JSON for you:

```
$ http \
  -X POST \
  https://httpbin.org/post \
  companyId=1234 \
  'name=John Doe'
```

Which generates the following request:

```HTTP
POST https://httpbin.org/post
Content-Type: application/json

{"companyId":"1234","name":"John Doe"}
```
