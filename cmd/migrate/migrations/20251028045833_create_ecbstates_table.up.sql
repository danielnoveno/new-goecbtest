/*
    file:           cmd/migrate/migrations/20251028045833_create_ecbstates_table.up.sql
    description:    Migration introducing the ECB states table with day-by-day status tracking.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `ecbstates` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tgl` DATE NOT NULL,
    `readstate` VARCHAR (20) DEFAULT "" NOT NULL,
    `ecbstate` VARCHAR (20) DEFAULT "" NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
