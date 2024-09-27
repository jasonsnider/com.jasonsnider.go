# com.jasonsnider.go

A go based implementation of jasonsnider.com, for the sake of learning Go.

# Getting stated

## Local Development

Clone the project 

```sh
git clone git@github.com:jasonsnider/com.jasonsnider.go.git
cd com.jasonsnider.go
```

Start the containers in the required environmnet production|development
```sh
docker compose --profile development up -d
docker compose --profile production up -d
```

Stop the containers
```sh
docker compose --profile development down --remove-orphans
docker compose --profile production down --remove-orphans
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


## Production Launch
- Login into the host machine and clone the project
- `cd com.jasonsnider.go`
- `cp .env.dist .env`
- `vim .env`
- Load the SSL certs into the private directory.
- `docker compose --profile production up -d`

- Test the DB connection
- Load the default data.

 ### private directory

```
- private
    - ssl
        - ${NGINX_HOST}
            - fullchain.pem
            - privkey.pem
```