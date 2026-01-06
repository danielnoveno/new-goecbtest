/*
    file:           cmd/migrate/migrations/20251028062256_create_compressors_table.up.sql
    description:    Migration creating the compressors reference table with metadata fields.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `compressors` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ctype` VARCHAR (4) NOT NULL DEFAULT "",
    `merk` VARCHAR (20) NOT NULL DEFAULT "",
    `type` VARCHAR (20) NOT NULL DEFAULT "",
    `itemcode` VARCHAR (20) NOT NULL DEFAULT "",
    `force_scan` INT UNSIGNED NOT NULL DEFAULT 1,
    `familycode` VARCHAR (20) NOT NULL DEFAULT "",
    `status` VARCHAR (20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
