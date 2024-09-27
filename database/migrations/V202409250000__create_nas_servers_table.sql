DROP TABLE IF EXISTS `nas_servers`;
CREATE TABLE `nas_servers` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `nas_id` integer NOT NULL,
  `name` varchar(200),
  `root_path` varchar(200),
  `internal_path` varchar(200),
  `created_at` timestamp NOT NULL DEFAULT (now())
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
