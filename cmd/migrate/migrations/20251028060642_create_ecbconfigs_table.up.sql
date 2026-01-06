/*
    file:           cmd/migrate/migrations/20251028060642_create_ecbconfigs_table.up.sql
    description:    Migration adding configurable key-value storage for ECB settings.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `ecbconfigs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `section` VARCHAR (20) DEFAULT "" NOT NULL,
    `variable` VARCHAR (20) DEFAULT "" NOT NULL,
    `value` TEXT DEFAULT "" NOT NULL,
    `ordering` VARCHAR(20) DEFAULT "000" NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
