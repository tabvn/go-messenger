sed -i '' 's/messenger?charset=utf8mb4/messengerdev?charset=utf8mb4/g' ./config/config.go;
sed -i '' 's/3007/3008/g' ./config/config.go;
sed -i '' 's/var\/www\/messenger/var\/www\/messenger-dev/g' ./config/config.go;
rm ./messenger
GOOS=linux GOARCH=amd64 go build
scp ./messenger ar:/var/www/messenger-dev/
scp -r ./public ar:/var/www/messenger-dev/
sed -i '' 's/messengerdev?charset=utf8mb4/messenger?charset=utf8mb4/g' ./config/config.go;
sed -i '' 's/3008/3007/g' ./config/config.go;
sed -i '' 's/var\/www\/messenger-dev/var\/www\/messenger/g' ./config/config.go;