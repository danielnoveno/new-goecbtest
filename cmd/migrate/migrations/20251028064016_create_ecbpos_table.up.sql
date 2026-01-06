/*
    file:           cmd/migrate/migrations/20251028064016_create_ecbpos_table.up.sql
    description:    Migration adding the ECB purchase order tracking table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `ecbpos` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `workcenter` VARCHAR(20) NOT NULL,
    `po` VARCHAR(20) NOT NULL,
    `sn` VARCHAR(20) NOT NULL,
    `ctype` INT UNSIGNED NOT NULL,
    `updated_by` INT UNSIGNED NOT NULL,
    `status` VARCHAR(20) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
