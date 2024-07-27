DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
	`id` bigint AUTO_INCREMENT PRIMARY KEY,
	`session_id` char(36) UNIQUE NOT NULL,
	`user_id` char(36) NOT NULL,
	`refresh_token` varchar(255) UNIQUE NOT NULL,
	`created_at` timestamp NOT NULL DEFAULT (now()),
	`expired_at` timestamp NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `sessions` ADD CONSTRAINT `fk_sessions_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);