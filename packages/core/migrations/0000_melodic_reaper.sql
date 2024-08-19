CREATE TABLE `card` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`user_id` char(30) NOT NULL,
	`stripe_payment_method_id` varchar(255) NOT NULL,
	`brand` text NOT NULL,
	`expiration_month` int NOT NULL,
	`expiration_year` int NOT NULL,
	`last4` char(4) NOT NULL,
	CONSTRAINT `card_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique` UNIQUE(`user_id`,`stripe_payment_method_id`)
);
--> statement-breakpoint
CREATE TABLE `cart_item` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`user_id` char(30) NOT NULL,
	`product_variant_id` char(30) NOT NULL,
	`quantity` int NOT NULL,
	CONSTRAINT `cart_item_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique` UNIQUE(`user_id`,`product_variant_id`)
);
--> statement-breakpoint
CREATE TABLE `cart` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`user_id` char(30) NOT NULL,
	`shipping_id` char(30),
	`card_id` char(30),
	CONSTRAINT `cart_id` PRIMARY KEY(`id`),
	CONSTRAINT `cart_user_id_unique` UNIQUE(`user_id`)
);
--> statement-breakpoint
CREATE TABLE `order_item` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`order_id` char(30) NOT NULL,
	`product_variant_id` char(30),
	`description` text,
	`quantity` int NOT NULL,
	`amount` bigint NOT NULL,
	CONSTRAINT `order_item_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique` UNIQUE(`order_id`,`product_variant_id`)
);
--> statement-breakpoint
CREATE TABLE `order` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`stripe_payment_intent_id` text NOT NULL,
	`user_id` char(30) NOT NULL,
	`shipping_address` json NOT NULL,
	`shipping_amount` bigint NOT NULL,
	`card` json NOT NULL,
	`tracking_number` text,
	`tracking_url` text,
	`label_url` text,
	`shippo_order_id` text,
	`shippo_label_id` text,
	`time_printed` timestamp(3),
	CONSTRAINT `order_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `product` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`name` varchar(255) NOT NULL,
	`description` text NOT NULL,
	CONSTRAINT `product_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `product_variant` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`product_id` char(30) NOT NULL,
	`name` varchar(255) NOT NULL,
	`price` bigint NOT NULL,
	CONSTRAINT `product_variant_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `user_fingerprint` (
	`user_id` char(30) NOT NULL,
	`fingerprint` varchar(255) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	CONSTRAINT `primary` PRIMARY KEY(`user_id`,`fingerprint`)
);
--> statement-breakpoint
CREATE TABLE `user_shipping` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`user_id` char(30) NOT NULL,
	`address` json NOT NULL,
	CONSTRAINT `user_shipping_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `user` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`name` varchar(255),
	`email` varchar(255),
	`fingerprint` varchar(255),
	`stripe_customer_id` varchar(255) NOT NULL,
	`flags` json DEFAULT ('{}'),
	CONSTRAINT `user_id` PRIMARY KEY(`id`),
	CONSTRAINT `user_fingerprint_unique` UNIQUE(`fingerprint`)
);
--> statement-breakpoint
ALTER TABLE `card` ADD CONSTRAINT `card_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `cart_item` ADD CONSTRAINT `cart_item_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `cart_item` ADD CONSTRAINT `cart_item_product_variant_id_product_variant_id_fk` FOREIGN KEY (`product_variant_id`) REFERENCES `product_variant`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `cart` ADD CONSTRAINT `cart_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `cart` ADD CONSTRAINT `cart_shipping_id_user_shipping_id_fk` FOREIGN KEY (`shipping_id`) REFERENCES `user_shipping`(`id`) ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `cart` ADD CONSTRAINT `cart_card_id_card_id_fk` FOREIGN KEY (`card_id`) REFERENCES `card`(`id`) ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `order_item` ADD CONSTRAINT `order_item_order_id_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `order`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `order_item` ADD CONSTRAINT `order_item_product_variant_id_product_variant_id_fk` FOREIGN KEY (`product_variant_id`) REFERENCES `product_variant`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `order` ADD CONSTRAINT `order_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `product_variant` ADD CONSTRAINT `product_variant_product_id_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `user_fingerprint` ADD CONSTRAINT `user_fingerprint_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `user_shipping` ADD CONSTRAINT `user_shipping_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE cascade ON UPDATE no action;