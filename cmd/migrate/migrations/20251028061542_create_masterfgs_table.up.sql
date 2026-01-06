/*
    file:           cmd/migrate/migrations/20251028061542_create_masterfgs_table.up.sql
    description:    Migration creating master finished goods metadata table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `masterfgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `mattype` VARCHAR (10) NOT NULL,
    `matdesc` VARCHAR (50) NOT NULL,
    `fgtype` VARCHAR (20) NOT NULL,
    `aging_tipes_id` INT NOT NULL DEFAULT 0,
    `kdbar` VARCHAR (20) NOT NULL,
    `warna` VARCHAR (20) NOT NULL,
    `lotinv` VARCHAR (20) NOT NULL,
    `attrib` VARCHAR (100) NOT NULL,
    `category` VARCHAR (20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
