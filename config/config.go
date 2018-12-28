package config

var WorkingDir = "/var/www/messenger"
var (
	Production      = false
	UploadDir       = WorkingDir + "/storage"
	PublicDir       = WorkingDir + "/public"
	MysqlConnectUrl = "messenger:messenger@tcp(127.0.0.1:3306)/messenger?charset=utf8mb4&collation=utf8mb4_unicode_ci"
	Port            = 3007
	PrivateAvatar = "/sites/default/files/default_images/anonymous.jpg"
)
