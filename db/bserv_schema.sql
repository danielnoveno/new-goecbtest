/*
   file:           db/bserv_schema.sql
   description:    Schema snapshot untuk DB BSERV (tabel ecb & cab_master)
*/

CREATE TABLE IF NOT EXISTS `ecb` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tgl` DATE NOT NULL,
    `jam` TIME NOT NULL,
    `gbj` VARCHAR(25) NOT NULL,
    `spc` VARCHAR(20) NOT NULL,
    `line` INT UNSIGNED NOT NULL,
    `ctype` VARCHAR(20) NOT NULL,
    `mfgpoststs` VARCHAR(20) NOT NULL DEFAULT '',
    `cdesc` VARCHAR(100) NOT NULL DEFAULT '',
    `ccode` VARCHAR(30) NOT NULL DEFAULT '',
    `tipe` VARCHAR(20) NOT NULL DEFAULT '',
    `tipemfg` VARCHAR(20) NOT NULL DEFAULT '',
    `lotinv` VARCHAR(20) NOT NULL DEFAULT '',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `ecb_tgl_jam_gbj_unique` (`tgl`, `jam`, `gbj`),
    KEY `ecb_gbj_spc_index` (`gbj`, `spc`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cab_master` (
    `cab_id` VARCHAR(20) NOT NULL,
    `desc1` VARCHAR(100) NOT NULL DEFAULT '',
    `desc2` VARCHAR(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`cab_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Seed contoh untuk simulasi remote duplicate-check
INSERT INTO `cab_master` (cab_id, desc1, desc2) VALUES
  ('SPCREF', 'REF DOOR', 'GASKET'),
  ('SPCPWM', 'PWM TUB', 'PULSATOR'),
  ('SPCVAC', 'VAC FOAM', 'SIDEWALL');

INSERT INTO `ecb` (tgl, jam, gbj, spc, line, ctype, mfgpoststs, cdesc, ccode, tipe, tipemfg, lotinv, created_at, updated_at) VALUES
  (CURDATE(), '07:55:00', 'RF1098000001', 'SPCREFOLD01', 1, 'RFA1', 'ECB_BL', 'Panasonic XA-10', 'CMPA10', 'REF', 'REF', 'LOT-REF-OLD', NOW(), NOW()),
  (CURDATE(), '08:10:00', 'RF1298000002', 'SPCREFOLD02', 2, 'RFB1', 'ECB_OK', 'LG BL-22', 'CMPB22', 'REF', 'REF', 'LOT-REF-B', NOW(), NOW()),
  (CURDATE(), '09:05:00', 'PW1098000003', 'SPCPWMOLD03', 3, 'PWM1', 'ECB_OK', 'Toshiba TD-12', 'CMPD12', 'PWM', 'PWM', 'LOT-PWM-A', NOW(), NOW());
