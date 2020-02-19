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
-- Table structure for table `approve_records`
--

DROP TABLE IF EXISTS `approve_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `approve_records` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `github` varchar(1023) COLLATE utf8_bin NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approve_records`
--

LOCK TABLES `approve_records` WRITE;
/*!40000 ALTER TABLE `approve_records` DISABLE KEYS */;
/*!40000 ALTER TABLE `approve_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auto_merge_white_names`
--

DROP TABLE IF EXISTS `auto_merge_white_names`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auto_merge_white_names` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `username` varchar(1023) COLLATE utf8_bin NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auto_merge_white_names`
--

LOCK TABLES `auto_merge_white_names` WRITE;
/*!40000 ALTER TABLE `auto_merge_white_names` DISABLE KEYS */;
/*!40000 ALTER TABLE `auto_merge_white_names` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auto_merges`
--

DROP TABLE IF EXISTS `auto_merges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auto_merges` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pr_id` int(11) DEFAULT NULL,
  `owner` varchar(1023) COLLATE utf8_bin DEFAULT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `started` tinyint(1) NOT NULL,
  `status` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auto_merges`
--

LOCK TABLES `auto_merges` WRITE;
/*!40000 ALTER TABLE `auto_merges` DISABLE KEYS */;
/*!40000 ALTER TABLE `auto_merges` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `black_names`
--

DROP TABLE IF EXISTS `black_names`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `black_names` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `username` varchar(1023) COLLATE utf8_bin NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `black_names`
--

LOCK TABLES `black_names` WRITE;
/*!40000 ALTER TABLE `black_names` DISABLE KEYS */;
/*!40000 ALTER TABLE `black_names` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cherry_prs`
--

DROP TABLE IF EXISTS `cherry_prs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cherry_prs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pr_id` int(11) DEFAULT '0',
  `from_pr` int(11) NOT NULL DEFAULT '0',
  `owner` varchar(255) COLLATE utf8_bin DEFAULT '',
  `repo` varchar(255) COLLATE utf8_bin DEFAULT '',
  `title` text COLLATE utf8_bin,
  `head` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `base` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `body` text COLLATE utf8_bin NOT NULL,
  `created_by_bot` tinyint(1) NOT NULL DEFAULT '0',
  `try_time` int(11) NOT NULL DEFAULT '0',
  `success` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `pr` (`pr_id`,`owner`,`repo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cherry_prs`
--

LOCK TABLES `cherry_prs` WRITE;
/*!40000 ALTER TABLE `cherry_prs` DISABLE KEYS */;
/*!40000 ALTER TABLE `cherry_prs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `issue_redelivers`
--

DROP TABLE IF EXISTS `issue_redelivers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `issue_redelivers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `issue_id` int(11) DEFAULT NULL,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `channel` varchar(1023) COLLATE utf8_bin NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `issue_redelivers`
--

LOCK TABLES `issue_redelivers` WRITE;
/*!40000 ALTER TABLE `issue_redelivers` DISABLE KEYS */;
/*!40000 ALTER TABLE `issue_redelivers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `label_checks`
--

DROP TABLE IF EXISTS `label_checks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `label_checks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pr_id` int(11) NOT NULL DEFAULT '0',
  `owner` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `repo` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `title` text COLLATE utf8_bin,
  `has_label` tinyint(1) DEFAULT '0',
  `send_notice` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `label_checks`
--

LOCK TABLES `label_checks` WRITE;
/*!40000 ALTER TABLE `label_checks` DISABLE KEYS */;
/*!40000 ALTER TABLE `label_checks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pull_requests`
--

DROP TABLE IF EXISTS `pull_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pull_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pr_id` int(11) NOT NULL DEFAULT '0',
  `owner` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `repo` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `title` text COLLATE utf8_bin,
  `label` text COLLATE utf8_bin,
  `merge` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `pr` (`owner`,`pr_id`,`repo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pull_requests`
--

LOCK TABLES `pull_requests` WRITE;
/*!40000 ALTER TABLE `pull_requests` DISABLE KEYS */;
/*!40000 ALTER TABLE `pull_requests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pull_status_checks`
--

DROP TABLE IF EXISTS `pull_status_checks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pull_status_checks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pull_id` int(11) DEFAULT NULL,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `label` varchar(1023) COLLATE utf8_bin NOT NULL,
  `duration` int(11) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pull_status_checks`
--

LOCK TABLES `pull_status_checks` WRITE;
/*!40000 ALTER TABLE `pull_status_checks` DISABLE KEYS */;
/*!40000 ALTER TABLE `pull_status_checks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pull_status_controls`
--

DROP TABLE IF EXISTS `pull_status_controls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pull_status_controls` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pull_id` int(11) DEFAULT NULL,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `label` varchar(1023) COLLATE utf8_bin NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pull_status_controls`
--

LOCK TABLES `pull_status_controls` WRITE;
/*!40000 ALTER TABLE `pull_status_controls` DISABLE KEYS */;
/*!40000 ALTER TABLE `pull_status_controls` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `slack_users`
--

DROP TABLE IF EXISTS `slack_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slack_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `github` varchar(1023) COLLATE utf8_bin NOT NULL DEFAULT '',
  `email` varchar(1023) COLLATE utf8_bin NOT NULL DEFAULT '',
  `slack` varchar(1023) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `slack_users`
--

LOCK TABLES `slack_users` WRITE;
/*!40000 ALTER TABLE `slack_users` DISABLE KEYS */;
/*!40000 ALTER TABLE `slack_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `white_names`
--

DROP TABLE IF EXISTS `white_names`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `white_names` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(1023) COLLATE utf8_bin NOT NULL,
  `repo` varchar(1023) COLLATE utf8_bin NOT NULL,
  `username` varchar(1023) COLLATE utf8_bin NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `white_names`
--

LOCK TABLES `white_names` WRITE;
/*!40000 ALTER TABLE `white_names` DISABLE KEYS */;
/*!40000 ALTER TABLE `white_names` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-02-09  6:48:14
