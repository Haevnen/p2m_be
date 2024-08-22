DROP TABLE IF EXISTS `links`;
CREATE TABLE `links` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `ticket_id` bigint NOT NULL,
  `link` VARCHAR(2083) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `ticket_id` bigint NOT NULL,
  `user_id` char(36) NOT NULL,
  `comment` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `histories`;
CREATE TABLE `histories` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `ticket_id` bigint NOT NULL,
  `action` text,
  `performed_by` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `links` ADD FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);
ALTER TABLE `comments` ADD FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);
ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);
ALTER TABLE `histories` ADD FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);
ALTER TABLE `histories` ADD FOREIGN KEY (`performed_by`) REFERENCES `users` (`user_id`);