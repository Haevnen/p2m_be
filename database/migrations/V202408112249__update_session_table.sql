ALTER TABLE `sessions`
MODIFY COLUMN `refresh_token` VARCHAR(512) UNIQUE NOT NULL;

CREATE INDEX idx_user_id ON `sessions` (`user_id`);
CREATE INDEX idx_session_id ON `sessions` (`session_id`);