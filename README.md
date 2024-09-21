# com.jasonsnider.com

A go based implementation of jasonsnider.com, for the sake of learning Go.

## Commands

Start the server
```sh
go run server.go -mode=server
```

Hash a password
```sh
go run server.go -mode=hash -password="<password>"
```

Compare plain-text to a hash
```sh
go run server.go -mode=check -password="<password>" -hashvalue="<hash>"
```