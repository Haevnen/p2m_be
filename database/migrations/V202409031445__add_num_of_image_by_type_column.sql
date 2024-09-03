ALTER TABLE `tickets` 
    ADD COLUMN `num_of_single_image` INT NOT NULL DEFAULT 0 AFTER `created_by`,
    ADD COLUMN `num_of_multiple_image` INT NOT NULL DEFAULT 0 AFTER `num_of_single_image`;