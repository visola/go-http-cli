### Authentication

You can also configure authentication from your profile. Basic and bearer are supported. An example of basic
authentication would look like the following:

```yaml
auth:
  type: basic
  username: myUsername
  password: myPassword
```

Your username and password will be automatically encoded accordingly to [RFC 7617](https://tools.ietf.org/html/rfc7617) and set as a header for your requests:

<pre>
$ http +httpbin /get

GET https://httpbin.org/get
<strong>Authorization: Basic bXlVc2VybmFtZTpteVBhc3N3b3Jk</strong>
Content-type: application/json
...
</pre>

Bearer (also known as token) authentication is also supported:

```yaml
auth:
  type: bearer
  token: myVerySecureToken
```

<pre>
$ http +httpbin /get

GET https://httpbin.org/get
<strong>Authorization: Bearer myVerySecureToken</strong>
Content-type: application/json
...
</pre>
