# go-graphql
Create GraphQL Api Service + Realtime PubSub use Golang

## Installation

* MYSQL (currently use MYSQL database)
* Change database connection in /config/config.go
* Import schema.sql to database
```
mysql -u root -p YOUR-DATABASE-NAME < schema.sql
```

## start Server

```
go run main.go

```

## React client

https://github.com/tabvn/messenger_client

<img src="https://firebasestorage.googleapis.com/v0/b/tabvn-fireshot.appspot.com/o/shots%2FQrC4k82w1uVqSO8ckTnvisBko7l1%2F-LIVDnqNVwxN4hWma2MU.png?alt=media&token=f2ab391b-a23c-47f3-9d9c-1a860e11559f" />
<img src="https://firebasestorage.googleapis.com/v0/b/tabvn-fireshot.appspot.com/o/shots%2FQrC4k82w1uVqSO8ckTnvisBko7l1%2F-LKHhNv8a-lSH7YOhewq.png?alt=media&token=707998bf-a48e-40ca-b50e-debbc4cd20fc" />

## Integration
* This chat service can be integrated into your website whatever language/CMS you are using.
