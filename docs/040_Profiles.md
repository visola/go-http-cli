## Profiles

`go-http-cli` can use profile files which are just YAML files in a special location.
The location by default points to `${user.home}/go-http-cli` but it can be configured through the
environment variable `GO_HTTP_PROFILES`.

**IMPORTANT! Please don't store your passwords on plain text files! Use this only for local/development environments.**

To activate a profile just add `+profileName` as part of your arguments. In this case, it would look for a
`${user.home}/go-http-cli/profileName.{yml|yaml}` file. It will fail if it can't find it.

What can you do with profiles? Many things:

### Base URL

You can set a base URL for all your calls into that profile. A simple example:

```yaml
baseURL:
  https://httpbin.org/
```

or 

```yaml
baseURL: https://httpbin.org/
```

Then calling `/ip` would automatically add the base URL like the following:

```
$ http +httpbin /ip

GET https://httpbin.org/ip
...
```

The path can be a relative path or an absolute path. The algorithm is very simple, it just concatenates
`baseURL` with `URL` making sure only one `/` will exist between the two. You can also override a `baseURL`
by passing a full URL from the command line (starting with `http` or `https`).

### Headers

Setting up headers for all your requests is a pain. So you can put it in your profile like the following:

```yaml
headers:
  Content-type:
    - application/json
```

or 

```yaml
headers:
  Content-type: application/json
```

And the headers will be added automatically:

```
http +httpbin /ip

GET https://httpbin.org/ip
Content-type: application/json
...
```
