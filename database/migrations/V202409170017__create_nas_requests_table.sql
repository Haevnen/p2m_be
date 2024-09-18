DROP TABLE IF EXISTS `nas_requests`;
CREATE TABLE `nas_requests` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `nas_id` integer NOT NULL,
  `payload` text,
  `status` ENUM ('FAILED', 'DONE') NOT NULL,
  `error` text,
  `created_at` timestamp NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;