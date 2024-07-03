CREATE TABLE IF NOT EXISTS tasks (
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),
  `linkedProjectID` CHAR(36) NOT NULL,
  `description` TEXT NOT NULL,
  `completed` BOOLEAN NOT NULL DEFAULT FALSE,

  PRIMARY KEY (id),
  FOREIGN KEY (linkedProjectID) REFERENCES projects(id)
);