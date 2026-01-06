/*
    file:           cmd/migrate/migrations/20261028062754_create_ecbdatas_table.up.sql
    description:    Migration creating the ecbdatas production log table with default values.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE TABLE IF NOT EXISTS `ecbdatas` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tgl` DATE NOT NULL DEFAULT "1920-01-01",
    `jam` TIME NOT NULL DEFAULT "00:00:00",
    `wc` VARCHAR (20) NOT NULL DEFAULT "",
    `prdline` VARCHAR (20) NOT NULL DEFAULT "",
    `ctgr` VARCHAR (255) NOT NULL DEFAULT "",
    `sn` VARCHAR (25) NOT NULL DEFAULT "",
    `fgtype` VARCHAR (20) NOT NULL DEFAULT "",
    `spc` VARCHAR (20) NOT NULL DEFAULT "",
    `comptype` VARCHAR (20) NOT NULL DEFAULT "",
    `compcode` VARCHAR (30) NOT NULL DEFAULT "",
    `po` VARCHAR (20) NOT NULL DEFAULT "",
    `status` VARCHAR (20) NOT NULL DEFAULT "",
    `sendsts` VARCHAR (20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
