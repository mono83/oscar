# Oscar

[![GitHub (pre-)release](https://img.shields.io/github/release/mono83/oscar/all.svg)](https://github.com/mono83/oscar/releases)
[![Docker Automated buil](https://img.shields.io/docker/automated/mono83/oscar.svg)](https://hub.docker.com/r/mono83/oscar/)

Toolset to run functional tests, written in [Lua](https://www.lua.org/).


## Oscar key features

* It is quick
* Imperative, not declarative syntax. This means you can use control structures, loops and custom functions
* Test cases are written in [Lua](https://www.lua.org/) - universal language
* Test cases, written using Oscar, are simple and can be easily refactored and generates human-readable diffs
  for code review 

## Quickstart

Download Oscar

```
go get -i https://github.com/mono83/oscar
```

Make simple test case and save it as `simple.lua`

```lua
-- Import Oscar module
local o = require("oscar")

-- Register test case
o.add(
    "Simple test case", -- Test case name
    function (tc)       -- Test case body. First and the only argument - TestCase object
        tc:log("This is simple test case body")
        tc:assertEquals("foo", "foo")
    end
)
```

Run it

```
oscar run simple.lua
```

## Command line arguments

### `oscar run <file.lua>`

* `--verbose|-v` Verbose mode - display detailed information about invocation, assertions and errors
* `--quiet|-q` Quiet mode - disables any output
* `--env|-e` INI file with values for global environment values map
* `--filter|-f` Filters test cases by name using regular expression
* `--lib|-l` Loads library lua file

## Test cases metadata and dependencies

There is possibility to provide additional metadata to test case:

```lua
-- Import Oscar module
local o = require("oscar")

-- Register test case
o.add(
    name, -- Test case name
    callback, -- Test case body
    metadata -- Additional metadata, optional
)
```

Metadata is simple Lua table with next keys:

| Key | Description | Example |
| --- | ----------- | ------- |
| `depends` | Describes test case dependency (by name). If that test case fails, current wont be invoked | `{depends="Simple GET with JSON"}` |


## Variables and placeholders

Placeholders have syntax `${name}` and are automatically replaced in most method:

```lua
tc:log("Running under ${lua.engine}")
```

Oscar uses variables map to replace placeholders in own methods. There are two major map classes:

* Global environment map. Values in this map are configured during start up (using `-e` option and 
  ini configuration file) and are immutable during test process.
* Test case values map. Values belongs to own test case and can be modified using `:set` method at
  runtime. Furthermore, all Oscar data providers (like `:httpPost` method) writes response and other
  data in this values map. 
  
Test case is allowed to read only own values map, but if value for requested key is not found, it
automatically falls back to global environment map. Reading can be achieved using `:get` method
or during variable interpolation process.

## Module-level functions 

### `init(function)`

Registers initialization (setUp) functions, that will be invoked only once for whole file. This 
function can contain heavy reusable calculations and environment establishment.

### `add(name, function)`

Registers new test case with name, passes as first argument. Second arguments stands for callback function,
that will be invoked on test case execution. Upon invocation, `TestCase` object will be passed to callback
as first and only argument

## RSA object

RSA is helper object, used to work with RSA public/private keys and certificates.

### Constructor

Use `RSA.new(len)` to generate public/private key pairs with desired `len` length and TTL one hour.

```lua
local r = RSA.new(2048)
```

### Certificate export

Syntax: `:exportCertificate64`

Outputs certificate in PEM base64 format. 

### Signatures

Syntax: `:signSHA256String64(data)`

Signs provided `data` using private key and `SHA256` hashing and returns base64-encoded value.

Syntax: `:signSHA256Int64(int...)`

Packs multiple int64 values (any amount) into one single BigEndian byte array, signs it and 
then encodes result using Base64.


## TestCase object

Each test case callback will automatically receive `TestCase` object as argument with following methods.

### Logging

Syntax: `:log(message)`

Outputs arbitrary message.

SyntaxL `:info(message)`

Outputs arbitrary message but with contrast color.

### Variables 

Syntax: `:get(name)`

Returns variable, identified by `name`

Syntax: `:set(name, value)`

Sets new `value` for variable, identified by `name`

Syntax: `:export(name, value)`

Sets new `value` for variable, identified by `name`. This value is set on top-level scope and can be exported 
to ini file.


### Assertions

Syntax: `:assertEquals(actual, expected [, description])`

Aliases: `:assert`, `:assertEq`

Performs equality check. All arguments casted to strings.

Syntax: `:assertJSONXPath(xpath, expected [, description])`

Parses last `http.response.body` as JSON, finds value under XPath expression and performs equality check for it.


### HTTP Requests

Syntax: `:httpGet(url)`

Performs HTTP request using `GET` method and writes response data into variables.

| Variable name | Meaning |
| ------------- | --------|
| `http.elapsed` | Time (in milliseconds), taken by request|
| `http.response.code` | HTTP status code |
| `http.response.length` | Response body length, in bytes |
| `http.response.body` | Full response body |
| `http.response.header.<name>` | Multiple values. Each response header will have own key |
| `http.request.url` | HTTP request URL |
| `http.request.length` | Request body length, in bytes |
| `http.request.body` | Full request body |
| `http.response.request.<name>` | Multiple values. Each request header will have own key |



Syntax: `:httpPost(url, body)`

Performs HTTP request using `POST` method and writes response data into variables.

| Variable name | Meaning |
| ------------- | --------| 
| `http.elapsed` | Time (in milliseconds), taken by request|
| `http.response.code` | HTTP status code |
| `http.response.length` | Response body length, in bytes |
| `http.response.body` | Full response body |
| `http.response.header.<name>` | Multiple values. Each response header will have own key |
| `http.request.url` | HTTP request URL |
| `http.request.length` | Request body length, in bytes |
| `http.request.body` | Full request body |
| `http.response.request.<name>` | Multiple values. Each request header will have own key |

### JSON 

Syntax: `:jsonXPath(path, body)`
         
Invokes JSON XPath query from `path` on `body` and returns invocation result. Interpolation also works
 
 ```lua
 local v = tc.jsonXPath("$.foo.bar", '{"foo":{"bar": 10}}')
 -- v = "10"
 ``` 
 
 ### Time
 
 Syntax: `:unix()`
 
 Returns current Unix timestamp, in seconds
 
 Syntax: `:sleep(millis)`
 
 Pauses invocation for requested amount of milliseconds
 
 ### Codecs
 
 ## Base64
 
 Syntax: `:stringToBase64(value)`
 
 Encodes provided string value as Base64
 
 Syntax: `:packInt64ToBase64(int64...)`
 
 Packs multiple int64 values (any amount) into one single BigEndian byte array and then encodes 
 it using Base64.
