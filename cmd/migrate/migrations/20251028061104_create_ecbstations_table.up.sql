/*
    file:           cmd/migrate/migrations/20251028061104_create_ecbstations_table.up.sql
    description:    Migration creating the ECB stations table to store station configuration.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `ecbstations` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ipaddress` VARCHAR (20) NOT NULL,
    `location` VARCHAR (50) NOT NULL,
    `mode` VARCHAR (20) NOT NULL,
    `linetype` VARCHAR (20) NOT NULL,
    `lineids` TEXT NOT NULL,
    `lineactive` INT NOT NULL,
    `ecbstate` VARCHAR (20) NOT NULL,
    `theme` VARCHAR (20) NOT NULL,
    `tacktime` INT NOT NULL,
    `workcenters` TEXT NOT NULL,
    `status` VARCHAR (20) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
