DROP TABLE IF EXISTS `clients`;
CREATE TABLE `clients` (
	`id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` VARCHAR(20) NOT NULL,
	`editing_style` TEXT,
	`requirements` TEXT,
	`others` TEXT,
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
	`created_at` TIMESTAMP NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;