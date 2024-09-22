# com.jasonsnider.com

A go based implementation of jasonsnider.com, for the sake of learning Go.

# Getting stated

## Local Development

Clone the project 

```sh
git clone git@github.com:jasonsnider/com.jasonsnider.go.git
cd com.jasonsnider.go
```

Start the containers
```sh
docker compose up -d
```

Stop the containers
```sh
docker compose down --remove-orphans
```

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