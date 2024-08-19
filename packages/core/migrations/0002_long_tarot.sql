ALTER TABLE `cart` DROP FOREIGN KEY `cart_shipping_id_user_shipping_id_fk`;
--> statement-breakpoint
ALTER TABLE `cart` ADD CONSTRAINT `cart_shipping_id_user_shipping_id_fk` FOREIGN KEY (`shipping_id`) REFERENCES `user_shipping`(`id`) ON DELETE set null ON UPDATE no action;