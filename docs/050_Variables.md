### Variables

Variables can be added to the profile or passed from the command line to either override from the
profile or add new values.

Variables will be replaced in the URL, body or headers. Just use them as `{variableName}` and it will
be automatically replaced.

In your profile:

```yaml
variables:
  companyId: 123456
```

Then you can make a call like the following:

```
$ http +httpbin /ip?companyId={companyId}

GET https://httpbin.org/ip?companyId=123456
Content-type: application/json
...
```
