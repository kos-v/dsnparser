## DSN Parser

This library parses a DSN of the form:  
![](doc/structure.png)  
Each of the sections is optional, so, for example, a DSN can be without schema or credentials.  

##### Using
```go
import "github.com/kos-v/dsnparser"

dsn := dsnparser.Parse("mysql://user:password@example.local:3306/dbname?tblsprefix=fs_");
dsn.GetScheme()            // string "mysql"
dsn.GetUser()              // string "user"
dsn.GetPassword()          // string "password"
dsn.GetHost()              // string "example.local"
dsn.GetPort()              // string "3306"
dsn.GetPath()              // string "dbname"
dsn.GetParam("tblsprefix") // string "fs_"
dsn.GetParams()            // map[string]string [tblsprefix: fs_]
dsn.HasParam("tblsprefix") // bool true
dsn.GetRaw()               // string "mysql://user:password@example.local:3306/dbname?tblsprefix=fs_"
```

##### Examples:
- `user:password@example.local`
- `example.local:65535`
- `user:password@`
- `socket:///foo/bar.sock`

##### Credentials
`user:password@` - user and password.  
`user@` or `user:@` - only user, without password.  
`:password@` - only password, without user.  

##### Escaping
You can escape the ":" and "@" characters in credentials, as well as the "=" and "&" characters in the extra options. To do this, specify "\\" before the desired character.  
Examples:  
`us\:e\@r:p\@ssw\:ord@` -> `us:e@r:p@ssw:ord@`  
`us\:e\@r:p\@ssw\:ord@` -> `us:e@r:p@ssw:ord@`  
`?key1=foo \& bar&key2=baz \= quux` -> `?key1=foo & bar&key2=baz = quux`  
`?foo\&key=fool val&bar\=key=bar val` -> `?foo&key=fool val&bar=key=bar val`  
