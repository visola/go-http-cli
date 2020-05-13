# go-http-cli
![Build](https://github.com/visola/go-http-cli/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/visola/go-http-cli)](https://goreportcard.com/report/github.com/visola/go-http-cli)
[![Maintainability](https://api.codeclimate.com/v1/badges/dda852b53b0e76299f8c/maintainability)](https://codeclimate.com/github/visola/go-http-cli/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/dda852b53b0e76299f8c/test_coverage)](https://codeclimate.com/github/visola/go-http-cli/test_coverage)

An HTTP client inspired by [Curl](https://github.com/curl/curl), [Postman](https://www.getpostman.com/) and [httpie](https://github.com/jakubroztocil/httpie) made with :heart: in Go.


## Table of Content

[Header](#header)<br />[Getting Started](#getting-started)<br />[Basic Usage](#basic-usage)<br />[Profiles](#profiles)<br />[Variables](#variables)<br />[Named Requests](#named-requests)<br />[Authentication](#authentication)<br />[Building from source](#building-from-source)<br />

## Getting Started

Download the latest [release](https://github.com/visola/go-http-cli/releases) for your platform.

Unzip it and put all three binaries in you path adding something like this to your `~/.bash_profile`:

```bash
export PATH=$PATH:/path/to/extracted/root
```

There's also an auto-completion helper for bash. You can add the following to your `~/.bash_profile`:

```bash
complete -f -C go-http-completion http
```

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

## Building from source

This project uses [Gradle](https://www.gradle.org) as the build system. To run it, the only thing
you'll need is to have Java installed. If you have a Go workspace [correctly setup](https://golang.org/doc/code.html)
you can run the following to run a complete build.

```bash
cd $GOPATH
mkdir -p src/github.com/visola
cd src/github.com/visola
git clone https://github.com/visola/go-http-cli.git
cd go-http-cli
./gradlew updateDependencies updateLinter build
```

There are also tasks to build specific packages for the `amd64` arch. You can build all the packages
using:

```bash
./gradlew buildPackages
```

or for one specific platform:

```bash
./gradlew packageDarwinAMD64
```

