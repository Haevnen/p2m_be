DROP TABLE IF EXISTS `users`;
CREATE TABLE users (
                             id                      BIGINT(20) AUTO_INCREMENT PRIMARY KEY,
                             nick_name                    VARCHAR(20) NOT NULL,
                             email                    VARCHAR(200) NOT NULL,
                             password              VARCHAR(100) NULL,
                             type               ENUM('FULLTIME', 'FREELANCER') NOT NULL,
                            is_active bit
);
