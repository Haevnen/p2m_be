ALTER TABLE `users` ADD COLUMN `is_unassigned` TINYINT(1) NOT NULL DEFAULT 0;

INSERT INTO `users` (user_id, nick_name, email, password_hashed, contract_type, is_active, is_admin, is_unassigned) values ("dc59bc05-8146-460e-a2b1-1ba5479cc63a", 'unassigned', 'unassigned@gmail.com', '$2a$10$3WF8TBrT1n8KsOTC.lA7xukLjvD3YFR1nx4OJgDnOH36pDMFtlvS2', 'FULLTIME', 1, 0, 1);