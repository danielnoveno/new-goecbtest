/*
    file:           cmd/migrate/migrations/20251028042239_create_navigations_table.up.sql
    description:    Migration creating the navigations menu table with hierarchy metadata.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS navigations (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `parent_id` INT UNSIGNED DEFAULT NULL,
    `icon` VARCHAR(100) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255) DEFAULT NULL,
    `url` VARCHAR(255) NOT NULL,
    `route` VARCHAR(255) DEFAULT NULL,
    `mode` INT UNSIGNED NOT NULL DEFAULT 1,
    `urutan` INT UNSIGNED NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
