ALTER TABLE `users` ADD COLUMN `is_unassigned` TINYINT(1) NOT NULL DEFAULT 0;

INSERT INTO `users` (user_id, nick_name, email, password_hashed, contract_type, is_active, is_admin, is_unassigned) values ("1a629d6e-6441-4a6a-b5ac-2487e3a569ed", 'unassigned', 'unassigned@gmail.com', '$2a$10$3WF8TBrT1n8KsOTC.lA7xukLjvD3YFR1nx4OJgDnOH36pDMFtlvS2', 'FULLTIME', 1, 0, 1);