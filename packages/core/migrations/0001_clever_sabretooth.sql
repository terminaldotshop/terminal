ALTER TABLE `cart_item` DROP FOREIGN KEY `cart_item_user_id_user_id_fk`;
--> statement-breakpoint
ALTER TABLE `cart_item` ADD CONSTRAINT `cart_item_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;