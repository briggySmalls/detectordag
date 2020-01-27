#
# Create tables
#

# User accounts
CREATE TABLE `accounts` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(128) NOT NULL,
  `password` VARCHAR(128) NOT NULL,
  PRIMARY KEY (id)
);

# Devices
CREATE TABLE `devices` (
  `id` BINARY(16),
  `account_id` INT(11),
  PRIMARY KEY (id),
  FOREIGN KEY (account_id) REFERENCES accounts(id)
);

# Emails
CREATE TABLE `emails` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(128) NOT NULL,
  `account_id` INT(11) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (account_id) REFERENCES accounts(id)
);
