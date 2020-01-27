# Seed database with some dummy values

# Create my account
INSERT INTO `accounts` (`username`, `password`) VALUES ('briggySmalls', 'password');
SET @accountId := (SELECT `id` FROM `accounts` WHERE `username` = 'briggySmalls');

# Create a device (my mac)
INSERT INTO `devices` (`id`, `account_id`)
    VALUES (UUID_TO_BIN('D856C48C-C006-53EE-B8AE-270DC96EB1F6', true), @accountId);

# Create some contact emails
INSERT INTO `emails` (`email`, `account_id`)
    VALUES ('wasso14@hotmail.com', @accountId);
INSERT INTO `emails` (`email`, `account_id`)
    VALUES ('briggysSmalls90@gmail.com', @accountId);
