package config

var (
	Production      = false
	UploadDir       = "/var/www/messenger/storage"
	PublicDir       = "/var/www/messenger/public"
	MysqlConnectUrl = "messenger:messenger@tcp(127.0.0.1:3306)/messenger?charset=utf8mb4&collation=utf8mb4_unicode_ci"
)
