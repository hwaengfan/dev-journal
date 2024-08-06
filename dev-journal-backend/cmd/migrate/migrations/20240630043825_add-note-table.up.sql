CREATE TABLE IF NOT EXISTS notes (
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),
  `userID` CHAR(36) NOT NULL,
  `linkedProjectID` CHAR(36) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `favorited` CHAR(5) NOT NULL DEFAULT "False",
  `tags` JSON NOT NULL,
  `dateCreated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastEdited` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  FOREIGN KEY (userID) REFERENCES users(id),
  FOREIGN KEY (linkedProjectID) REFERENCES projects(id)
);