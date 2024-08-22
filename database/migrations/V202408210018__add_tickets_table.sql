DROP TABLE IF EXISTS `tickets`;
CREATE TABLE `tickets` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `title` varchar(255) NOT NULL,
  `status` ENUM ('BACKLOG', 'IN_PROGRESS', 'READY_TO_QC', 'QC_VERIFYING', 'QC_DONE', 'DONE') NOT NULL,
  `qc_id` char(36) NOT NULL,
  `editor_id` char(36) NOT NULL,
  `priority` ENUM ('NORMAL', 'HIGH') NOT NULL DEFAULT 'NORMAL',
  `client_id` integer NOT NULL,
  `description` text,
  `created_by` ENUM ('AUTO', 'MANUAL') NOT NULL DEFAULT 'AUTO',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `tickets` ADD CONSTRAINT `fk_tickets_editor_users` FOREIGN KEY (`editor_id`) REFERENCES `users` (`user_id`);
ALTER TABLE `tickets` ADD CONSTRAINT `fk_tickets_qc_users` FOREIGN KEY (`qc_id`) REFERENCES `users` (`user_id`);
ALTER TABLE `tickets` ADD CONSTRAINT `fk_tickets_clients` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`);