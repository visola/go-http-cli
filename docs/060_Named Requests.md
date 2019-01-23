### Named Requests

You can preconfigure requests inside a profile and then call them by name using `@requestName`. For example,
if you had this in your profile:

```yaml
requests:
  makePost:
    url: /post
    body: '{
      "username": "{username}",
      "companyId": {companyId}
    }'
```

You could make the following call:

```
$ http +httpbin @makePost -V username=test

POST https://httpbin.org/post
Content-type: application/json
>> { "username": "test", "companyId": 123456 }
...
```
