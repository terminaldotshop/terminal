ALTER TABLE `order` MODIFY COLUMN `user_id` char(30);--> statement-breakpoint
ALTER TABLE `order` MODIFY COLUMN `card` json;