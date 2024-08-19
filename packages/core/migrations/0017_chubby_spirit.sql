CREATE TABLE `product_variant_inventory` (
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`product_variant_id` char(30) NOT NULL,
	`inventory_id` char(30) NOT NULL,
	CONSTRAINT `product_variant_inventory_product_variant_id_inventory_id_pk` PRIMARY KEY(`product_variant_id`,`inventory_id`)
);
--> statement-breakpoint
ALTER TABLE `product_variant_inventory` ADD CONSTRAINT `product_variant_inventory_product_variant_id_fk` FOREIGN KEY (`product_variant_id`) REFERENCES `product_variant`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `product_variant_inventory` ADD CONSTRAINT `product_variant_inventory_inventory_inventory_id_fk` FOREIGN KEY (`inventory_id`) REFERENCES `inventory`(`id`) ON DELETE cascade ON UPDATE no action;
