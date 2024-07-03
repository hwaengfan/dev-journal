CREATE TABLE IF NOT EXISTS users (
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),
  `firstName` VARCHAR(255) NOT NULL,
  `lastName` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,

  PRIMARY KEY (id),
  UNIQUE KEY (email)
);