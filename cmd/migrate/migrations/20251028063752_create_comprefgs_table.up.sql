/*
    file:           cmd/migrate/migrations/20251028063752_create_comprefgs_table.up.sql
    description:    Migration creating compressor reference finished goods mapping table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `comprefgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ctype` VARCHAR(4) NOT NULL DEFAULT "",
    `barcode` VARCHAR(20) NOT NULL DEFAULT "",
    `status` VARCHAR(20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
