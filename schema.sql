-- MySQL dump 10.13  Distrib 5.7.25, for macos10.14 (x86_64)
--
-- Host: localhost    Database: messenger
-- ------------------------------------------------------
-- Server version	5.7.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `archived`
--

DROP TABLE IF EXISTS `archived`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `archived` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`group_id`,`user_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `archived_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE,
  CONSTRAINT `archived_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archived`
--

LOCK TABLES `archived` WRITE;
/*!40000 ALTER TABLE `archived` DISABLE KEYS */;
INSERT INTO `archived` VALUES (1,65,10),(2,65,19);
/*!40000 ALTER TABLE `archived` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attachments`
--

DROP TABLE IF EXISTS `attachments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `attachments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `message_id` int(11) unsigned DEFAULT '0',
  `file_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `message_id` (`message_id`),
  KEY `file_id` (`file_id`),
  CONSTRAINT `attachments_ibfk_1` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE,
  CONSTRAINT `attachments_ibfk_2` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attachments`
--

LOCK TABLES `attachments` WRITE;
/*!40000 ALTER TABLE `attachments` DISABLE KEYS */;
INSERT INTO `attachments` VALUES (1,595,127),(2,595,128),(3,595,129);
/*!40000 ALTER TABLE `attachments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `blocked`
--

DROP TABLE IF EXISTS `blocked`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blocked` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `author` int(11) unsigned DEFAULT NULL,
  `user` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`author`,`user`),
  KEY `user` (`user`),
  CONSTRAINT `blocked_ibfk_1` FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `blocked_ibfk_2` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `blocked`
--

LOCK TABLES `blocked` WRITE;
/*!40000 ALTER TABLE `blocked` DISABLE KEYS */;
/*!40000 ALTER TABLE `blocked` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deleted`
--

DROP TABLE IF EXISTS `deleted`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deleted` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL,
  `message_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`user_id`,`message_id`),
  KEY `message_id` (`message_id`),
  CONSTRAINT `deleted_ibfk_2` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE,
  CONSTRAINT `deleted_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deleted`
--

LOCK TABLES `deleted` WRITE;
/*!40000 ALTER TABLE `deleted` DISABLE KEYS */;
/*!40000 ALTER TABLE `deleted` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `files` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `original` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `type` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `size` int(11) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `files_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=130 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `files`
--

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;
INSERT INTO `files` VALUES (41,10,'dc7bb5ee95b5b70791f60c49.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534588700),(42,10,'179316dd63dfbd3347c09937.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534588718),(43,10,'6a22f1b090dd5cb66eda7e97.jpg','newyork2.jpg','image/jpeg',508854,1534588738),(44,10,'9d83314cc8a455b2b09539f8.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534588828),(45,10,'fde29452be8341203b957592.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534588843),(46,10,'be6fa503052eb59d5c6b36d7.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534597668),(47,10,'1891570d5e8d6f581bac3924.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534598613),(48,10,'75feecd85011e244ec524314.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534598744),(49,10,'27919b4c6dda79bf06f81b96.pdf','INVOICE_10_Tabvn.pdf','application/pdf',21828,1534666456),(50,10,'3e814fa51e599f5cb81fb991.pdf','INVOICE_10_Tabvn.pdf','application/pdf',21828,1534666470),(51,10,'8a8f5577a96148a47e3b255b.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534666536),(52,10,'c174e865657f80ef7f9f17d3.pdf','screencapture-dev-clearinsightbenefits-2018-07-27-16_39_01.pdf','application/pdf',1047555,1534668235),(53,10,'96f7c67ed036fafed5b3334c.pdf','Starting Out with C++ from Control Structures to Objects.pdf','application/pdf',7332501,1534668235),(54,10,'03fe90b030877f2f3171d172.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534668790),(55,10,'723f8a124b3dab8e64cacf48.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534668846),(56,10,'e953925cc4147dc39bc23f28.pdf','NIFTIT.Toan.Non.EcommerceStream.06.12.2018.pdf','application/pdf',305994,1534670658),(57,10,'57aa911c576010ec8bdfe3c3.pdf','RTL 1.10 UX Test.pdf','application/pdf',33066,1534670658),(58,10,'2923dc230c6987155d0dcb12.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534670793),(59,10,'ac8896be0ce9bfd4bd522f03.txt','huuduc (1).pub','text/plain; charset=utf-8',403,1534670793),(60,10,'a73da5163d58e11f1dddfa1d.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534670961),(61,10,'027544d9ebce4860f0312d6d.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671008),(62,10,'84ccabdf82398006f5a2afb1.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671055),(63,10,'8e866db3ab22e0da3f25fede.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671107),(64,10,'bb64bf285114a78839348cca.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671149),(65,10,'fe0a3a7f40981af0968bf312.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671244),(66,10,'5f8d04af29497a1e7a9ff628.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534671382),(67,10,'518f73cc5d6f3cc959a60038.pdf','screencapture-dev-clearinsightbenefits-2018-07-27-16_39_01.pdf','application/pdf',1047555,1534671454),(68,10,'9756da86a0fbcda63ba18cd6.pdf','Starting Out with C++ from Control Structures to Objects.pdf','application/pdf',7332501,1534671454),(69,10,'6a49bfc185395e9eae30321c.jpg','newyork2.jpg','image/jpeg',508854,1534671486),(70,10,'6f1925402da4f99e573a8b88.jpg','newyork2.jpg','image/jpeg',508854,1534671754),(71,10,'845095417ae523e154305445.jpg','newyork2.jpg','image/jpeg',508854,1534671837),(72,10,'4bf2a948f8d883b85a7c3c53.jpg','newyork2.jpg','image/jpeg',508854,1534671864),(73,10,'175cf2dacabad6397019fb7f.jpg','newyork2.jpg','image/jpeg',508854,1534671870),(74,10,'f0f20761e23751d4aa203f76.jpg','newyork2.jpg','image/jpeg',508854,1534671918),(75,10,'e23029560ce9faeef080e6a8.jpg','newyork2.jpg','image/jpeg',508854,1534671926),(76,10,'261c6157f88050b9af6a31a1.jpg','newyork2.jpg','image/jpeg',508854,1534672142),(77,10,'6499871aa3b94fd180f099e0.jpg','newyork2.jpg','image/jpeg',508854,1534672157),(78,10,'e6eaa1836c4b09eb039b078c.jpg','newyork1.jpg','image/jpeg',1161745,1534673046),(79,10,'3d113091dd5a881019a3dd00.jpg','newyork2.jpg','image/jpeg',508854,1534673046),(80,10,'c1495976f2d2467e74f9d307.jpg','tyler_1.jpeg','image/jpeg',46209,1534673073),(81,10,'a303ef21a60ab3ffea97aaf6.jpg','tyler.jpeg','image/jpeg',163981,1534673073),(82,10,'3d1f1fa98b4b43ce83182117.bin','ar-messenger-sidebar-collapsed.psd','application/octet-stream',8150544,1534673107),(83,10,'016dcced14b996f437e4ac98.jpg','newyork1.jpg','image/jpeg',1161745,1534673210),(84,10,'43db58903502b802e01c3eb7.jpg','newyork2.jpg','image/jpeg',508854,1534673210),(85,10,'d6666063a6002d338fb003de.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534673290),(86,10,'04cc46c11708edd2a5ae4e7e.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534673377),(87,10,'59ca28884b35bcf5cdfe1a18.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534673404),(88,10,'b54f5d4b4802a06384e2d3d1.bin','ca7b282f-126b-4c76-8700-80d67ce1af2e_tkbhkyimydtupdt.xls','application/octet-stream',2039808,1534673438),(89,10,'167eb4d9408771ba1cf71a66.pdf','JAVIS_HOME_Quotation_Dealer_20180808.pdf','application/pdf',429306,1534674574),(90,10,'56009dae47adac0dae611600.pdf','The C++Standard Library - 2nd Edition.pdf','application/pdf',14470662,1534674579),(91,10,'3d34a4ac565eddc6f55a728b.zip','Thiet_ke_nha_from_wonder.vn_CAD_DWG.zip','application/zip',2784268,1534674579),(92,10,'eebabb35962e8435e8b29d37.jpg','tyler_1.jpeg','image/jpeg',46209,1534752047),(93,10,'3c8ea5f3b9bd175eff6f0761.bin','ca7b282f-126b-4c76-8700-80d67ce1af2e_tkbhkyimydtupdt.xls','application/octet-stream',2039808,1534819951),(94,10,'aef88e1e21f54621a355713f.bin','CN_Congnghethongtin.xls','application/octet-stream',287744,1534819951),(95,10,'260c640486e535557fe2738e_cn-congnghethongtin.xls','CN_Congnghethongtin.xls','application/octet-stream',287744,1534821240),(96,10,'e88d8989cb7f5cc1211baec2_ar-messenger-modals1.psd','ar-messenger-modals(1).psd','application/octet-stream',18930894,1534821267),(97,10,'0124739d57100c5a4ea7547a_golangwebsocketpubsub-service.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534891206),(98,10,'e37e7b3d90744dd3ffe373fe_tyler-1.jpeg','tyler_1.jpeg','image/jpeg',46209,1534909688),(99,10,'f9016a1da0898d20c79cccfd_golangwebsocketpubsub-service.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534950746),(100,10,'aafac54feefac4b796d4b196_golangwebsocketpubsub-service.png','GoLangWebsocketPubSub Service.png','image/png',171706,1534950830),(101,10,'8b25e9ab64902a7f9bc02912_invoice-10-tabvn.pdf','INVOICE_10_Tabvn.pdf','application/pdf',21828,1534950864),(102,10,'faf2cc43fe90549c693899a2_ar-messenger.psd','ar-messenger.psd','application/octet-stream',14178364,1534950899),(103,10,'bf85fe46cb3f32ed635fd850_huuduc-1.pub','huuduc (1).pub','text/plain; charset=utf-8',403,1534950909),(104,10,'ff8fd60de4b46a697e545306_huuduc-1.pub','huuduc (1).pub','text/plain; charset=utf-8',403,1534950920),(105,10,'c252e92554b62d7a5c27ad75_newyork1.jpg','newyork1.jpg','image/jpeg',1161745,1535065504),(106,10,'cd933f4de2e663a53e7d2d47_tyler-1.jpeg','tyler_1.jpeg','image/jpeg',46209,1535065920),(107,10,'393854ddb1b5893a8a8d3fae_tyler-1.jpeg','tyler_1.jpeg','image/jpeg',46209,1535066074),(108,10,'72175cde1df481c86343480d_tyler-1.jpeg','tyler_1.jpeg','image/jpeg',46209,1535066132),(109,10,'5b4d3d601d68eec58ec6d15a_newyork1.jpg','newyork1.jpg','image/jpeg',1161745,1535066221),(110,10,'913e1a4410a7f62ef6f349d4_tyler.jpeg','tyler.jpeg','image/jpeg',163981,1535066404),(111,10,'08c3b6eda27f9282a3c1f3d7_invoice-10-tabvn.pdf','INVOICE_10_Tabvn.pdf','application/pdf',21828,1535066540),(112,10,'43129031c69b64d1b1efd226_tyler.jpeg','tyler.jpeg','image/jpeg',163981,1535066554),(113,10,'11034945631d6899ad965d6a_kadane-algorithm.png','kadane-Algorithm.png','image/png',18328,1544519921),(114,10,'bac3de0a6ad13e2cac01f782_kadane-algorithm.png','kadane-Algorithm.png','image/png',18328,1544519988),(115,10,'a11b3b328b476d9f48cb4ec3_abottadorable.png','abott@adorable.png','image/png',13578,1544578699),(116,10,'e37baec2d094db9f249dd65c_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544578751),(117,10,'6ace7d0450477a91321868f6_abottadorable.png','abott@adorable.png','image/png',13578,1544578884),(118,10,'304162c863deafc48040ff00_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544578920),(119,10,'dfb923158af12456ad486b3d_abottadorable.png','abott@adorable.png','image/png',13578,1544579000),(120,10,'5eec007aa06d54611e262c34_abottadorable.png','abott@adorable.png','image/png',13578,1544579018),(121,10,'8a9dcb694a304fa431154131_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544579057),(122,10,'94b3a5f93dba84590addb824_abottadorable.png','abott@adorable.png','image/png',13578,1544579130),(123,10,'19616eab2a97925492a7323c_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544579154),(124,10,'fc63877c2da4ba973b50528b_abottadorable.png','abott@adorable.png','image/png',13578,1544579191),(125,10,'985b2da1f86f884b84b03d31_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544579283),(126,10,'baa866973b59455e9be8e342_abottadorable.png','abott@adorable.png','image/png',13578,1544579297),(127,10,'991bf4bc29010d5e784344c2_47296421-1536513249782420-2214504344911347712-n.jpg','47296421_1536513249782420_2214504344911347712_n.jpg','image/jpeg',137595,1544751926),(128,10,'2f59e8cac22cb99f29360912_47400456-10210933851662881-7002594436941086720-o.jpg','47400456_10210933851662881_7002594436941086720_o.jpg','image/jpeg',142707,1544751926),(129,10,'8c1c135a24d91ea4b667e51d_47473832-10210933852982914-4844884806816759808-o.jpg','47473832_10210933852982914_4844884806816759808_o.jpg','image/jpeg',136417,1544751926);
/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friendship`
--

DROP TABLE IF EXISTS `friendship`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `friendship` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `friend_id` int(11) unsigned DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_friendship` (`user_id`,`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friendship`
--

LOCK TABLES `friendship` WRITE;
/*!40000 ALTER TABLE `friendship` DISABLE KEYS */;
INSERT INTO `friendship` VALUES (3,10,20,1,1545968206),(4,20,10,0,1545968206);
/*!40000 ALTER TABLE `friendship` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `groups`
--

DROP TABLE IF EXISTS `groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8 DEFAULT '',
  `user_id` int(11) unsigned DEFAULT '0',
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  FULLTEXT KEY `fulltext` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `groups`
--

LOCK TABLES `groups` WRITE;
/*!40000 ALTER TABLE `groups` DISABLE KEYS */;
INSERT INTO `groups` VALUES (62,'','bac3de0a6ad13e2cac01f782_kadane-algorithm.png',10,1535069458,1544519989),(63,'','',10,1535078477,1535078477),(64,'','',12,1535247198,1535247198),(65,'','',10,1535871830,1535871830),(66,'','',19,1543799941,1543799941),(68,'','',10,1543896626,1543896626),(69,'','',10,1543911735,1543911735),(70,'','',10,1543912027,1543912027),(71,'','',20,1543974200,1543974200),(72,'Test dtest ','baa866973b59455e9be8e342_abottadorable.png',20,1543976715,1544579298),(73,'','',10,1545976320,1545976320);
/*!40000 ALTER TABLE `groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `members`
--

DROP TABLE IF EXISTS `members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `members` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `group_id` int(11) unsigned DEFAULT '0',
  `blocked` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT '0',
  `accepted` tinyint(1) DEFAULT '0',
  `added_by` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_index` (`user_id`,`group_id`),
  KEY `group_id` (`group_id`),
  KEY `added_by` (`added_by`),
  CONSTRAINT `members_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `members_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE,
  CONSTRAINT `members_ibfk_3` FOREIGN KEY (`added_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=299 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `members`
--

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
INSERT INTO `members` VALUES (267,10,62,0,1535069458,1,NULL),(268,16,63,0,1535078477,1,NULL),(270,10,63,0,1535078477,1,NULL),(271,10,64,0,1535247198,1,NULL),(275,16,65,0,1535871830,1,NULL),(276,10,65,0,1535871830,1,NULL),(278,10,66,0,1543799941,1,NULL),(279,19,66,0,1543799941,1,NULL),(280,19,65,0,1543891474,1,NULL),(288,20,70,0,1543912027,1,10),(292,10,72,0,1543976715,1,20),(293,20,72,0,1543976715,1,20),(294,20,62,0,1544516943,0,NULL),(295,19,62,0,1544516947,0,NULL),(296,19,72,0,1544577285,1,NULL),(297,20,73,0,1545976320,0,10),(298,10,73,0,1545976320,1,10);
/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `messages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `group_id` int(11) unsigned DEFAULT '0',
  `body` text COLLATE utf8mb4_bin,
  `emoji` tinyint(1) DEFAULT '0',
  `gif` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `group_id` (`group_id`),
  FULLTEXT KEY `body` (`body`),
  CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=627 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (367,10,62,'f',0,'',1535069458,1535069458),(369,10,62,'fdsafa',0,'',1535070462,1535070462),(373,10,62,'fdsfasf',0,'',1535080433,1535080433),(375,10,62,'fsafsaf',0,'',1535080434,1535080434),(381,10,62,'f',0,'',1535080436,1535080436),(382,10,62,'sf',0,'',1535080436,1535080436),(385,10,62,'s',0,'',1535080437,1535080437),(386,10,62,'af',0,'',1535080437,1535080437),(387,10,62,'saf',0,'',1535080437,1535080437),(402,10,64,'?',1,'',1535248620,1535248620),(403,10,64,'how are you',0,'',1535248623,1535248623),(404,10,64,'doing',0,'',1535248624,1535248624),(406,10,64,'',0,'l4FGA9IjmH4u1NZbq',1535249018,1535249018),(407,10,64,'?',1,'',1535249142,1535249142),(408,10,64,'how are you',0,'',1535249353,1535249353),(409,10,64,'doing',0,'',1535249361,1535249361),(410,10,64,'how are you today',0,'',1535249928,1535249928),(411,10,64,'how are you today i \'m doing well',0,'',1535249967,1535249967),(412,10,64,'hi',0,'',1535249970,1535249970),(413,10,64,'',0,'',1535249976,1535249976),(414,10,64,'',0,'',1535249988,1535249988),(415,10,64,'',0,'',1535250025,1535250025),(416,10,64,'',0,'',1535250027,1535250027),(417,10,64,'',0,'',1535250028,1535250028),(418,10,64,'',0,'',1535250033,1535250033),(419,10,64,'',0,'',1535250035,1535250035),(420,10,64,'f',0,'',1535250146,1535250146),(421,10,64,'f',0,'',1535250155,1535250155),(422,10,64,'? cool ?',0,'',1535251838,1535251988),(423,10,64,'?',1,'',1535251991,1535252116),(424,10,64,'?',1,'',1535252125,1535252125),(425,10,64,'?',1,'',1535252626,1535254043),(426,10,64,'how ',0,'',1535252630,1535252630),(427,10,64,'good',0,'',1535252740,1535252740),(428,10,64,'job',0,'',1535252745,1535252745),(429,10,64,'how ',0,'',1535252782,1535252782),(430,10,64,'how',0,'',1535252835,1535252835),(431,10,64,'are you',0,'',1535252838,1535252838),(432,10,64,'gfsadfsafff  fsafsdfsa fsda fsa  f safsafsaf fdsafa',0,'',1535252894,1535252899),(433,10,64,'how are you\nfdsaf',0,'',1535253553,1535253553),(434,10,64,'line 1\nline 2 ?',0,'',1535253989,1535254021),(436,10,64,'http://tabvn.com ?',0,'',1535255279,1535255307),(437,10,64,'good job',0,'',1535255321,1535255328),(438,10,64,' toan hi alert(1) Toan',0,'',1535257599,1535257729),(439,10,64,'Toan',0,'',1535257761,1535257761),(440,10,64,'hi',0,'',1535274727,1535274727),(441,10,64,'how',0,'',1535274730,1535274730),(442,10,64,'hi',0,'',1535274854,1535274854),(443,10,64,'hi',0,'',1535275939,1535275939),(444,10,64,'???',1,'',1535640729,1535640729),(445,10,64,'???',1,'',1535640740,1535640740),(446,10,64,'?',1,'',1535640742,1535640742),(447,10,65,'hi',0,'',1535871830,1535871830),(448,10,65,'how are you',0,'',1535871835,1535871835),(449,10,64,'hi',0,'',1536028127,1536028127),(452,10,62,'hey 1',0,'',1536113593,1538035514),(453,10,62,'',0,'kEbbp4YovnpJmwWyKe',1537579372,1537579372),(454,10,62,'',0,'5T0yHHWsC39rtI82SH',1537579639,1537579639),(455,10,62,'1',0,'',1538035216,1538035216),(456,10,62,'2',0,'',1538035216,1538035216),(457,10,62,'3',0,'',1538035217,1538035217),(458,10,62,'4 cool',0,'',1538035217,1538035554),(459,10,62,'5 edit',0,'',1538035218,1538035533),(460,10,62,'f',0,'',1538036521,1538036521),(461,10,62,'f',0,'',1538036522,1538036522),(462,10,62,'g',0,'',1538036523,1538036523),(463,10,62,'j',0,'',1538036525,1538036525),(464,10,62,'toan 2',0,'',1538036624,1538036627),(465,10,62,'toan 14',0,'',1538036636,1538036643),(468,10,62,'gdgs',0,'',1538036694,1538036694),(469,10,62,'toan',0,'',1538036734,1538036734),(470,10,62,'fsaf 2',0,'',1538036820,1538036828),(471,10,62,'toan',0,'',1538037167,1538037167),(472,10,62,'toan',0,'',1538037205,1538037205),(473,10,62,'test',0,'',1538037255,1538037255),(474,10,62,'f',0,'',1538037271,1538037271),(475,10,62,'f',0,'',1538037281,1538037281),(476,10,62,'fdfsa',0,'',1538037318,1538037318),(477,10,62,'hi',0,'',1538037344,1538037344),(478,10,62,'how are you 1',0,'',1538037346,1538037350),(479,10,62,'',0,'l9WzIeAx3bToV7IM5F',1538124670,1538124670),(480,10,62,'',0,'l9WzIeAx3bToV7IM5F',1538124678,1538124678),(481,10,62,'',0,'l9WzIeAx3bToV7IM5F',1538124723,1538124723),(482,10,62,'',0,'l4pTii07Gypi3GFPy',1538124727,1538124727),(491,10,64,'hi',0,'',1538269916,1538269916),(492,10,64,'this is test \nvery long message',0,'',1538534180,1538534180),(493,10,64,'hi',0,'',1538535742,1538535742),(494,10,64,'test',0,'',1539331068,1539331068),(495,10,64,'test',0,'',1540352073,1540352073),(496,10,64,'http://drupal.org',0,'',1540352681,1540352681),(497,10,64,'http://drupal.org',0,'',1540352707,1540352707),(498,10,64,'http://tabvn.com',0,'',1540352779,1540352779),(499,10,64,'http://apple.com',0,'',1540352800,1540352800),(500,10,64,'http://apple.com',0,'',1540352833,1540352833),(501,10,65,'hi there',0,'',1541730156,1541730156),(502,10,65,'hi',0,'',1541730174,1541730174),(503,19,65,'Tyler here',0,'',1541730196,1541730196),(504,19,65,'Tyler here',0,'',1541730252,1541730252),(505,19,65,'Tyler here',0,'',1541730870,1541730870),(506,10,65,'how are you',0,'',1541730995,1541730995),(507,10,65,'toan here',0,'',1541731002,1541731002),(508,19,65,'this is tyler',0,'',1541731019,1541731019),(509,16,65,'ok alex here',0,'',1541731214,1541731214),(510,10,65,'toan here',0,'',1541731231,1541731231),(511,16,65,'hi',0,'',1541732782,1541732782),(512,16,65,'hi',0,'',1541732796,1541732796),(513,16,65,'hi',0,'',1541732799,1541732799),(514,16,65,'how are you',0,'',1541732801,1541732801),(515,16,65,'hi',0,'',1541732843,1541732843),(516,16,65,'how are you',0,'',1541732845,1541732845),(517,16,65,'ho',0,'',1541732851,1541732851),(518,16,65,'hi',0,'',1541732955,1541732955),(519,16,65,'how ar you',0,'',1541732961,1541732961),(520,16,65,'toan here',0,'',1541733010,1541733010),(521,16,65,'goo djob',0,'',1541733014,1541733014),(522,16,65,'hi',0,'',1541733021,1541733021),(523,16,65,'how are you',0,'',1541733027,1541733027),(524,16,65,'it is working',0,'',1541733033,1541733033),(525,16,65,'fine',0,'',1541733042,1541733042),(526,10,65,'?',1,'',1541737520,1541737520),(527,10,65,'fdsafafafs',0,'',1541925974,1541925974),(528,10,65,'how are you',0,'',1541926093,1541926093),(529,10,65,'fdsaf\nfdsaf\nfdasf\nfdsaf\nfdsaf\nfdsaf\nfdsaf\nfdsafa',0,'',1541926887,1541926887),(530,10,65,'hi',0,'',1541927867,1541927867),(531,10,65,'how ar eyou',0,'',1543799609,1543799609),(532,10,65,'are you ok ?',0,'',1543799618,1543799618),(533,10,65,'how are you',0,'',1543799860,1543799860),(534,19,65,'good',0,'',1543799886,1543799886),(535,19,65,'u ?',0,'',1543799894,1543799894),(536,19,65,'hi',0,'',1543799922,1543799922),(537,19,65,'yhou can not nsed',0,'',1543799932,1543799932),(538,19,66,'hi',0,'',1543799941,1543799941),(539,19,66,'fsdfa',0,'',1543799946,1543799946),(540,19,66,'fsdaf',0,'',1543799948,1543799948),(541,19,66,'fsafsdfa',0,'',1543799954,1543799954),(542,10,65,'hi',0,'',1543800637,1543800637),(543,19,65,'tyler here',0,'',1543800660,1543800660),(544,10,65,'good',0,'',1543800760,1543800760),(545,19,66,'hi',0,'',1543800789,1543800789),(546,10,65,'f',0,'',1543800815,1543800815),(547,19,65,'g',0,'',1543800818,1543800818),(550,10,68,'ji',0,'',1543896626,1543896626),(551,10,68,'ff',0,'',1543911143,1543911143),(552,20,68,'hi',0,'',1543911223,1543911223),(553,20,68,'how are you',0,'',1543911235,1543911235),(554,10,68,'West here',0,'',1543911241,1543911241),(555,20,68,'f',0,'',1543911284,1543911284),(556,20,68,'good',0,'',1543911288,1543911288),(557,20,68,'good job',0,'',1543911295,1543911295),(558,20,68,'f',0,'',1543911316,1543911316),(559,20,68,'fsdfsaf',0,'',1543911318,1543911318),(560,10,68,'how',0,'',1543911326,1543911326),(561,10,68,'hey',0,'',1543911333,1543911333),(562,20,68,'good',0,'',1543911336,1543911336),(563,10,68,'drupal',0,'',1543911341,1543911341),(564,20,68,'ff',0,'',1543911374,1543911374),(566,20,68,'hey',0,'',1543911389,1543911389),(567,20,68,'g',0,'',1543911408,1543911408),(568,20,68,'how',0,'',1543911409,1543911409),(569,10,68,'arere',0,'',1543911412,1543911412),(570,10,69,'hi',0,'',1543911735,1543911735),(571,10,69,'f',0,'',1543911981,1543911981),(572,20,69,'hi',0,'',1543911984,1543911984),(573,10,70,'g',0,'',1543912027,1543912027),(574,20,70,'hi',0,'',1543912070,1543912070),(575,10,70,'ff',0,'',1543912075,1543912075),(576,20,70,'hi',0,'',1543912182,1543912182),(577,20,70,'f',0,'',1543974174,1543974174),(578,20,70,'a',0,'',1543974178,1543974178),(579,20,71,'fsafsaf',0,'',1543974200,1543974200),(580,20,72,'hi',0,'',1543976715,1543976715),(581,10,72,'hi West',0,'',1543976734,1543976734),(582,20,72,'hi Toan',0,'',1543976738,1543976738),(583,20,72,'i \'m fine',0,'',1543976740,1543976740),(584,20,72,'hi',0,'',1543977210,1543977210),(585,20,72,'no sent',0,'',1543977225,1543977225),(586,20,72,'hi',0,'',1543977227,1543977227),(587,20,72,'how are yu',0,'',1543977229,1543977229),(588,10,72,'giid',0,'',1543977234,1543977234),(589,20,72,'hi',0,'',1543977238,1543977238),(590,20,72,'how ',0,'',1543977239,1543977239),(591,10,72,'?',1,'',1544491266,1544491266),(592,10,72,'?',1,'',1544491285,1544491285),(593,10,72,'? ? ? ',0,'',1544577215,1544577215),(594,10,72,'???',0,'',1544577226,1544577226),(595,10,72,'',0,'',1544751926,1544751926),(596,19,72,'hi',0,'',1545011362,1545011362),(597,19,72,'h',0,'',1545011377,1545011377),(598,19,72,'how are you',0,'',1545011380,1545011380),(599,19,72,'good',0,'',1545011394,1545011394),(600,19,66,'hi ',0,'',1545011519,1545011519),(601,19,66,'how are you',0,'',1545011523,1545011523),(602,19,66,'i \'m good',0,'',1545011539,1545011539),(603,19,66,'coo',0,'',1545011581,1545011581),(604,19,66,'good',0,'',1545011587,1545011587),(605,10,72,'how are yu',0,'',1545011600,1545011600),(606,10,72,'good',0,'',1545011621,1545011621),(607,10,72,'f',0,'',1545011626,1545011626),(608,19,72,'ok',0,'',1545011641,1545011641),(609,19,72,'thus is good',0,'',1545011647,1545011647),(610,19,72,'another test message',0,'',1545011666,1545011666),(611,19,72,'test',0,'',1545011673,1545011673),(612,19,72,'f',0,'',1545011686,1545011686),(613,19,72,'another test',0,'',1545011695,1545011695),(614,19,72,'tét',0,'',1545013337,1545013337),(615,19,72,'gôd',0,'',1545013340,1545013340),(616,10,62,'f',0,'',1545014029,1545014029),(617,10,73,'hi',0,'',1545976320,1545976320),(618,10,73,'hi',0,'',1545976358,1545976358),(619,10,73,'http://drupal.org',0,'',1546835792,1546835792),(620,10,73,'http://drupal.org',0,'',1546835976,1546835976),(621,10,73,'http://drupal.org',0,'',1546836000,1546836000),(622,10,73,'http://facebook.com',0,'',1546836135,1546836135),(623,10,73,'http://homekitvietnam.com/',0,'',1546836159,1546836159),(624,10,73,'http://homekitvietnam.com/',0,'',1546836179,1546836179),(625,10,73,'http://homekitvietnam.com',0,'',1546836721,1546836721),(626,19,66,'hi',0,'',1547116535,1547116535);
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `unread_insert_trigger` AFTER INSERT ON `messages` FOR EACH ROW begin
  insert into unreads (message_id, user_id)  SELECT new.id,members.user_id FROM members WHERE members.group_id = new.group_id AND members.accepted != 2 AND members.blocked = 0 AND members.user_id != new.user_id ;
  delete from archived where group_id = new.group_id ;
end */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `secrets`
--

DROP TABLE IF EXISTS `secrets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `secrets` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `secret` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `secrets`
--

LOCK TABLES `secrets` WRITE;
/*!40000 ALTER TABLE `secrets` DISABLE KEYS */;
INSERT INTO `secrets` VALUES (1,'8472e809-cda2-4b1e-8289-120373ca7f4b1');
/*!40000 ALTER TABLE `secrets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tokens`
--

DROP TABLE IF EXISTS `tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT NULL,
  `token` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tokens`
--

LOCK TABLES `tokens` WRITE;
/*!40000 ALTER TABLE `tokens` DISABLE KEYS */;
INSERT INTO `tokens` VALUES (1,10,'77a0199f-ee92-4766-9b50-54e5faa96834',1534291668),(18,18,'4a6a2d1e-bf08-4515-a8e9-91f781acb7d7',1535004015),(19,10,'ea7976b5-84b8-4ff4-a9dc-ebf923ee94f9',1535004267),(21,10,'8555a992-03f0-40ce-843a-f05fd5eaca23',1535021210),(29,10,'88498d31-9575-4248-b6a6-6fb1c9914c35',1541731148),(30,16,'eab81530-6792-49f4-a820-23819a57e2ac',1541731208),(31,19,'2b27079a-53e8-4420-b379-a741e867b378',1541731868),(32,20,'fdca132f-9516-43d8-b26d-d34fc2d178b2',1543911213);
/*!40000 ALTER TABLE `tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `unreads`
--

DROP TABLE IF EXISTS `unreads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `unreads` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0',
  `message_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `message_id` (`message_id`),
  CONSTRAINT `unreads_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `unreads_ibfk_2` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=943 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `unreads`
--

LOCK TABLES `unreads` WRITE;
/*!40000 ALTER TABLE `unreads` DISABLE KEYS */;
INSERT INTO `unreads` VALUES (808,16,526),(810,16,527),(813,16,528),(816,16,529),(819,16,530),(821,16,531),(824,16,532),(826,16,533),(829,16,534),(832,16,535),(835,16,536),(838,16,537),(845,16,542),(848,16,543),(851,16,544),(855,16,546),(858,16,547),(867,10,564),(869,10,566),(870,10,567),(871,10,568),(872,10,577),(873,10,578),(874,16,579),(876,20,581),(880,20,588),(881,20,591),(882,20,592),(883,20,593),(884,20,594),(885,20,595),(888,20,596),(891,20,597),(894,20,598),(897,20,599),(899,20,605),(902,20,606),(905,20,607),(909,20,608),(912,20,609),(915,20,610),(918,20,611),(921,20,612),(924,20,613),(927,20,614),(930,20,615),(932,20,616),(934,20,617),(935,20,618),(936,20,619),(937,20,620),(938,20,621),(939,20,622),(940,20,623),(941,20,624),(942,20,625);
/*!40000 ALTER TABLE `unreads` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT '0',
  `first_name` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `last_name` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `email` varchar(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `password` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `avatar` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `online` tinyint(1) DEFAULT '0',
  `custom_status` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `location` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `work` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `school` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `about` longtext CHARACTER SET utf8,
  `created` int(11) DEFAULT '0',
  `updated` int(11) DEFAULT '0',
  `published` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_email` (`email`),
  UNIQUE KEY `unique_uid` (`uid`),
  FULLTEXT KEY `fulltext` (`first_name`,`last_name`,`email`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (10,1,'admin','','toan@tabvn.com','$2a$10$XRj3pj8/lTgBD4hXlV9CkOo/G.ObfTXisAOtnwsPnAydGTEBohGMS','',0,'online',NULL,NULL,NULL,NULL,1534289600,1541731148,1),(16,4,'Alex','','alex@gmail.com','$2a$10$feLVXNrNx5lP37dDAG5f6OQSWcIjc/.XfWKAOn0jWFw3gOpoKC/1i','',0,NULL,NULL,NULL,NULL,NULL,1534998787,1541731208,1),(18,7,'Toan 7','Mr','toan7@tabvn.com','$2a$10$4yEwhoeC49PlVzUcsiw.8exkGmJ9Q0huyiJx3oXNJsjaQXfqtvBRS','',0,NULL,NULL,NULL,NULL,NULL,1535003566,1548899802,1),(19,2,'Tyler','','tyler@gmail.com','$2a$10$5z9tkPN/fSMFLxy3C.4RGus.ai7wU2vHkakQba94qsq5xdbTvRMs2','',0,NULL,NULL,NULL,NULL,NULL,1535021130,1541731868,1),(20,3,'West','','west@gmail.com','$2a$10$DRFshGRj7XRfgt90PzRE1.UQfNvR7Ez85.fyRFpq4o1Rv0C.oed/e','',0,NULL,NULL,NULL,NULL,NULL,1541731127,1543911213,2),(21,9,'Toan 7','Mr','toan99@tabvn.com','$2a$10$6FupurzhH4rKVWecDJ9I7ODvpsKMNQLMNHjabJZv7sJL6FypwE70G','',0,NULL,NULL,NULL,NULL,NULL,1548900077,1548900077,1);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-04-12  9:58:30
