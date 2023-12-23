CREATE TABLE IF NOT EXISTS `tmd_forzamotorsport2023` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `created_at` varchar(100) NOT NULL DEFAULT current_timestamp(),
   `IsRaceOn` tinyint(1) DEFAULT NULL,
   `TimestampMS` float DEFAULT NULL,
   `EngineMaxRpm` double DEFAULT NULL,
   `EngineIdleRpm` double DEFAULT NULL,
   `CurrentEngineRpm` decimal(10,0) DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;