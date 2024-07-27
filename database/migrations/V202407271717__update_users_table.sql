ALTER TABLE `users`
    CHANGE COLUMN `type` `contract_type` ENUM ('FULLTIME', 'FREELANCE') NOT NULL DEFAULT 'FULLTIME',
    CHANGE COLUMN `password` `password_hashed` varchar(255) NOT NULL,
    MODIFY COLUMN `id` int AUTO_INCREMENT,
    MODIFY COLUMN `nick_name` varchar(20) UNIQUE NOT NULL,
    MODIFY COLUMN `email` varchar(200) UNIQUE NOT NULL,
    MODIFY COLUMN `is_active` bit NOT NULL DEFAULT 1,
    ADD COLUMN `is_admin` bit NOT NULL DEFAULT 0,
    ADD COLUMN `user_id` char(36) UNIQUE NOT NULL AFTER `id`,
    ADD COLUMN `created_at` timestamp NOT NULL DEFAULT (now());

ALTER TABLE `users` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;