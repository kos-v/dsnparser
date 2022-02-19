## DSN Parser
[![Build Status](https://app.travis-ci.com/kos-v/dsnparser.svg?branch=main)](https://app.travis-ci.com/kos-v/dsnparser)
[![Coverage Status](https://codecov.io/gh/kos-v/dsnparser/branch/main/graph/badge.svg)](https://codecov.io/gh/kos-v/dsnparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/kos-v/dsnparser)](https://goreportcard.com/report/github.com/kos-v/dsnparser)

This library parses a DSN of the form:  
![](doc/structure.png)  
Each of the sections is optional, so, for example, a DSN can be without schema or credentials.  

##### Installation
```shell
go get github.com/kos-v/dsnparser
```

##### Using
```go
import "github.com/kos-v/dsnparser"

dsn := dsnparser.Parse("mysql://user:password@tcp(example.local:3306)/dbname?encoding=utf8mb4");
dsn.GetScheme()          // string "mysql"
dsn.GetUser()            // string "user"
dsn.GetPassword()        // string "password"
dsn.GetHost()            // string "example.local"
dsn.GetPort()            // string "3306"
dsn.GetPath()            // string "dbname"
dsn.GetParam("encoding") // string "utf8mb4"
dsn.GetParams()          // map[string]string [encoding: "utf8mb4"]
dsn.GetTransport()       // string "tcp"
dsn.HasParam("encoding") // bool true
dsn.GetRaw()             // string "mysql://user:password@example.local:3306/dbname?encoding=utf8mb4"
```

##### Examples:
- `user:password@example.local`
- `example.local:65535`
- `user:password@`
- `socket:///foo/bar.sock`
- `mysql://user:password@example.local/dbname`
- `mysql://example.local/?db=dbname&user=user&password=password`

##### Credentials
`user:password@` - user and password.  
`user@` or `user:@` - only user, without password.  
`:password@` - only password, without user.  

##### Escaping
You can escape the ":" and "@" characters in credentials, as well as the "=" and "&" characters in the extra options. To do this, specify "\\" before the desired character.  
Examples:  
`us\:e\@r:p\@ssw\:ord@` -> `us:e@r:p@ssw:ord`  
`us\:e\@r:p\@ssw\:ord@` -> `us:e@r:p@ssw:ord`  
`?key1=foo \& bar&key2=baz \= quux` -> `key1=foo & bar&key2=baz = quux`  
`?foo\&key=fool val&bar\=key=bar val` -> `foo&key=fool val&bar=key=bar val`  
