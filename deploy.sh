rm ./messenger
GOOS=linux GOARCH=amd64 go build
scp ./messenger ar:/var/www/messenger/
scp -r ./public ar:/var/www/messenger/
