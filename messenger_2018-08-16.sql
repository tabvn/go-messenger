# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.6.40)
# Database: messenger
# Generation Time: 2018-08-16 16:06:41 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table attachments
# ------------------------------------------------------------

DROP TABLE IF EXISTS `attachments`;

CREATE TABLE `attachments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `message_id` int(11) unsigned DEFAULT '0',
  `file_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `message_id` (`message_id`),
  KEY `file_id` (`file_id`),
  CONSTRAINT `attachments_ibfk_1` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE,
  CONSTRAINT `attachments_ibfk_2` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `attachments` WRITE;
/*!40000 ALTER TABLE `attachments` DISABLE KEYS */;

INSERT INTO `attachments` (`id`, `message_id`, `file_id`)
VALUES
	(1,5,1),
	(2,5,2);

/*!40000 ALTER TABLE `attachments` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table blocked
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocked`;

CREATE TABLE `blocked` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `author` int(11) unsigned DEFAULT NULL,
  `user` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`author`,`user`),
  KEY `user` (`user`),
  CONSTRAINT `blocked_ibfk_1` FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `blocked_ibfk_2` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table files
# ------------------------------------------------------------

DROP TABLE IF EXISTS `files`;

CREATE TABLE `files` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `name` varchar(255) DEFAULT NULL,
  `original` varchar(255) DEFAULT NULL,
  `type` varchar(50) DEFAULT NULL,
  `size` int(11) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `files_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;

INSERT INTO `files` (`id`, `user_id`, `name`, `original`, `type`, `size`, `created`)
VALUES
	(1,12,'tyler.jpeg','tyler_original.ong','image/png',100,0),
	(2,12,'tyler_1.jpeg','tyler_1_original.png','image/png',0,0),
	(4,12,'test.jpg','test_original.jpg','image/jpeg',100,0);

/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table friendship
# ------------------------------------------------------------

DROP TABLE IF EXISTS `friendship`;

CREATE TABLE `friendship` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `friend_id` int(11) unsigned DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_friendship` (`user_id`,`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `friendship` WRITE;
/*!40000 ALTER TABLE `friendship` DISABLE KEYS */;

INSERT INTO `friendship` (`id`, `user_id`, `friend_id`, `status`, `created`)
VALUES
	(1,10,12,1,1534303433),
	(2,12,10,1,1534303433),
	(4,10,14,1,1534310404),
	(5,14,10,1,1534310404);

/*!40000 ALTER TABLE `friendship` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table groups
# ------------------------------------------------------------

DROP TABLE IF EXISTS `groups`;

CREATE TABLE `groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT '',
  `user_id` int(11) unsigned DEFAULT '0',
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  FULLTEXT KEY `fulltext` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `groups` WRITE;
/*!40000 ALTER TABLE `groups` DISABLE KEYS */;

INSERT INTO `groups` (`id`, `title`, `avatar`, `user_id`, `created`, `updated`)
VALUES
	(13,'','',10,1534291571,1534291571),
	(21,'','',10,1534292167,1534292167),
	(22,'','',11,1534292194,1534292194),
	(23,'','',10,1534298771,1534298771);

/*!40000 ALTER TABLE `groups` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table members
# ------------------------------------------------------------

DROP TABLE IF EXISTS `members`;

CREATE TABLE `members` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `group_id` int(11) unsigned DEFAULT '0',
  `blocked` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_index` (`user_id`,`group_id`),
  KEY `group_id` (`group_id`),
  CONSTRAINT `members_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `members_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;

INSERT INTO `members` (`id`, `user_id`, `group_id`, `blocked`, `created`)
VALUES
	(19,10,13,0,1534291571),
	(20,11,13,0,1534291571),
	(42,10,21,0,1534292167),
	(43,11,21,0,1534292167),
	(44,12,21,0,1534292167),
	(45,11,22,0,1534292194),
	(46,12,22,0,1534292194),
	(47,10,23,0,1534298771),
	(48,12,23,0,1534298771);

/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table messages
# ------------------------------------------------------------

DROP TABLE IF EXISTS `messages`;

CREATE TABLE `messages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `group_id` int(11) unsigned DEFAULT '0',
  `body` text,
  `emoji` tinyint(1) DEFAULT '0',
  `gif` varchar(255) DEFAULT NULL,
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `group_id` (`group_id`),
  FULLTEXT KEY `body` (`body`),
  CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;

INSERT INTO `messages` (`id`, `user_id`, `group_id`, `body`, `emoji`, `gif`, `created`, `updated`)
VALUES
	(1,10,13,'Toan to Alex',0,'0',1534291571,1534291571),
	(2,10,21,'Toan to Alex and Tyler',0,'0',1534292167,1534292167),
	(3,11,22,'Alex to Tyler',0,'0',1534292194,1534292194),
	(4,10,23,'Hi Tyler how are you ?',0,'0',1534298771,1534298771),
	(5,12,23,'Tyler here how are you toan, attachments',0,'0',1534299553,1534299553),
	(6,10,23,'Toan here how are you Tyler',0,'0',1534409352,1534409352);

/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

DELIMITER ;;
/*!50003 SET SESSION SQL_MODE="NO_ENGINE_SUBSTITUTION" */;;
/*!50003 CREATE */ /*!50017 DEFINER=`root`@`localhost` */ /*!50003 TRIGGER `unread_insert_trigger` AFTER INSERT ON `messages` FOR EACH ROW begin
  insert into unreads (message_id, user_id)  SELECT new.id,members.user_id FROM members WHERE members.group_id = new.group_id AND members.blocked = 0 AND members.user_id != new.user_id ;
end */;;
DELIMITER ;
/*!50003 SET SESSION SQL_MODE=@OLD_SQL_MODE */;


# Dump of table secrets
# ------------------------------------------------------------

DROP TABLE IF EXISTS `secrets`;

CREATE TABLE `secrets` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `secret` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `secrets` WRITE;
/*!40000 ALTER TABLE `secrets` DISABLE KEYS */;

INSERT INTO `secrets` (`id`, `secret`)
VALUES
	(1,'8472e809-cda2-4b1e-8289-120373ca7f4b1');

/*!40000 ALTER TABLE `secrets` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tokens
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tokens`;

CREATE TABLE `tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT NULL,
  `token` varchar(255) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `tokens` WRITE;
/*!40000 ALTER TABLE `tokens` DISABLE KEYS */;

INSERT INTO `tokens` (`id`, `user_id`, `token`, `created`)
VALUES
	(1,10,'b7339bb5-a1e1-49ae-b932-7667e1d0db93',1534291668);

/*!40000 ALTER TABLE `tokens` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table unreads
# ------------------------------------------------------------

DROP TABLE IF EXISTS `unreads`;

CREATE TABLE `unreads` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `message_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `message_id` (`message_id`),
  CONSTRAINT `unreads_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `unreads_ibfk_2` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `unreads` WRITE;
/*!40000 ALTER TABLE `unreads` DISABLE KEYS */;

INSERT INTO `unreads` (`id`, `user_id`, `message_id`)
VALUES
	(1,11,1),
	(2,11,2),
	(3,12,2),
	(5,12,3),
	(6,12,4),
	(7,10,5),
	(8,12,6);

/*!40000 ALTER TABLE `unreads` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT '0',
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `email` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `avatar` varchar(255) DEFAULT NULL,
  `online` tinyint(1) DEFAULT '0',
  `custom_status` varchar(50) DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `work` varchar(255) DEFAULT NULL,
  `school` varchar(255) DEFAULT NULL,
  `about` longtext,
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_email` (`email`),
  UNIQUE KEY `unique_uid` (`uid`),
  FULLTEXT KEY `fulltext` (`first_name`,`last_name`,`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `uid`, `first_name`, `last_name`, `email`, `password`, `avatar`, `online`, `custom_status`, `location`, `work`, `school`, `about`, `created`, `updated`)
VALUES
	(10,1,'Toan','Nguyen','toan@tabvn.com','$2a$10$XRj3pj8/lTgBD4hXlV9CkOo/G.ObfTXisAOtnwsPnAydGTEBohGMS','https://api.adorable.io/avatars/100/abott@adorable.png',0,NULL,NULL,NULL,NULL,NULL,1534289600,1534291668),
	(11,2,'Alex','M','alex@tabvn.com','$2a$10$GjRosHggB8rYpMPX/iWG.OSxR6e9pNxGDeWG1sncEb4hJi4N3Mojy','',0,NULL,NULL,NULL,NULL,NULL,1534289616,1534289616),
	(12,3,'Tyler','C','tyler@tabvn.com','$2a$10$1bXQ3XIEmyN9Q6cA0qIHpe1NRNGwlsC4eHRkveJj44ufL1RH.s1Za','https://dev.addictionrecovery.com/sites/default/files/picture-38-1495649849_3.jpg',0,NULL,NULL,NULL,NULL,NULL,1534289632,1534289632),
	(14,4,'West','J','west@tabvn.com','$2a$10$rPNyk8yUVQktPRij9HKO0uNQ9igBXQQ758/Dx/iSUNVLTS1KVwhZu','',1,NULL,NULL,NULL,NULL,NULL,1534289664,1534289664);

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
