DROP TABLE IF EXISTS `session`;

CREATE TABLE IF NOT EXISTS `session` (
   `session_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'Session id',
   `contents` TEXT NOT NULL COMMENT 'Session data',
   `last_active` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Last active time',
   `expiration` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Expiration time',
   PRIMARY KEY (`session_id`),
   KEY `last_active` (`last_active`),
   KEY `expiration` (`expiration`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='session table';
