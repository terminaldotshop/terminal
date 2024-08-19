CREATE TABLE `inventory_record` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`inventory_id` char(30) NOT NULL,
	`quantity` int NOT NULL,
	`notes` text,
	CONSTRAINT `inventory_record_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `inventory` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`name` varchar(255) NOT NULL,
	`description` text,
	CONSTRAINT `inventory_id` PRIMARY KEY(`id`),
	CONSTRAINT `inventory_name_unique` UNIQUE(`name`)
);
--> statement-breakpoint
ALTER TABLE `inventory_record` ADD CONSTRAINT `inventory_record_inventory_id_inventory_id_fk` FOREIGN KEY (`inventory_id`) REFERENCES `inventory`(`id`) ON DELETE cascade ON UPDATE no action;