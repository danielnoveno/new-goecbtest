/*
    file:           cmd/migrate/migrations/20251028062102_create_mastersfgs_table.up.sql
    description:    Migration establishing the semi-finished goods reference table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `mastersfgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `plant` VARCHAR (10) NOT NULL,
    `mattype` VARCHAR (10) NOT NULL,
    `matdesc` VARCHAR (50) NOT NULL,
    `sfgtype` VARCHAR (20) NOT NULL,
    `sfgdesc` VARCHAR (50) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
