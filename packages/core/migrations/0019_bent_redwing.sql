CREATE TABLE `subscription` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`time_next` timestamp(3),
	`user_id` char(30) NOT NULL,
	`frequency` varchar(255) NOT NULL,
	`product_variant_id` char(30) NOT NULL,
	`quantity` int NOT NULL,
	`shipping_id` char(30) NOT NULL,
	`card_id` char(30) NOT NULL,
	CONSTRAINT `subscription_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique` UNIQUE(`user_id`,`product_variant_id`)
);
--> statement-breakpoint
ALTER TABLE `product` ADD `subscription` varchar(255);--> statement-breakpoint
ALTER TABLE `subscription` ADD CONSTRAINT `subscription_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `subscription` ADD CONSTRAINT `subscription_product_variant_id_product_variant_id_fk` FOREIGN KEY (`product_variant_id`) REFERENCES `product_variant`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `subscription` ADD CONSTRAINT `subscription_shipping_id_user_shipping_id_fk` FOREIGN KEY (`shipping_id`) REFERENCES `user_shipping`(`id`) ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `subscription` ADD CONSTRAINT `subscription_card_id_card_id_fk` FOREIGN KEY (`card_id`) REFERENCES `card`(`id`) ON DELETE no action ON UPDATE no action;
