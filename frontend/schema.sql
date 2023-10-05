CREATE TABLE `users` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(255) UNIQUE NOT NULL,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `role` ENUM ('root_admin', 'admin', 'user', 'subuser') NOT NULL DEFAULT "user",
  `permissions` TEXT NOT NULL,
  `use_totp` TINYINT(3) NOT NULL,
  `totp_secret` TEXT DEFAULT NULL,
  `totp_authenticated_at` TIMESTAMP DEFAULT NULL,
  `created_at` TIMESTAMP DEFAULT NULL,
  `updated_at` TIMESTAMP DEFAULT NULL
);

CREATE TABLE `activity_log` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `event` VARCHAR(255) NOT NULL,
  `ip` VARCHAR(45) NOT NULL,
  `description` TEXT,
  `actor_type` varchar(191) DEFAULT NULL,
  `actor_id` bigint(20) DEFAULT NULL,
  `timestamp` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE UNIQUE INDEX `unique_root_admin` ON `users` (`role`);

ALTER TABLE `activity_log` ADD FOREIGN KEY (`actor_id`) REFERENCES `users` (`id`);