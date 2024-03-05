CREATE TABLE IF NOT EXISTS tmd_forzamotorsport2023_bestlaps (
    id bigint(20) unsigned auto_increment NOT NULL,
    created_at varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'current_timestamp()' NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    CarOrdinal INT UNSIGNED DEFAULT NULL NULL,
    TrackOrdinal INT UNSIGNED DEFAULT NULL NULL,
    BestLap float DEFAULT NULL NULL,
    CarClass INT UNSIGNED DEFAULT NULL NULL,
    CarPerformanceIndex INT UNSIGNED DEFAULT NULL NULL,
    DrivetrainType INT UNSIGNED DEFAULT NULL NULL,
    NumCylinders INT UNSIGNED DEFAULT NULL NULL,
    Fuel float DEFAULT NULL NULL,
    LapNumber INT DEFAULT NULL NULL,
    RacePosition INT DEFAULT NULL NULL,
    CONSTRAINT `PRIMARY` PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
