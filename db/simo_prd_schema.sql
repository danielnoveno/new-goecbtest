/*
   file:           db/simo_prd_schema.sql
   description:    Schema snapshot untuk DB SIMO_PRD (mirror dari Laravel)
*/

-- Stations
CREATE TABLE IF NOT EXISTS `ecbstations` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ipaddress` VARCHAR(20) NOT NULL,
    `location` VARCHAR(50) NOT NULL,
    `mode` VARCHAR(20) NOT NULL,
    `linetype` VARCHAR(20) NOT NULL,
    `lineids` TEXT NOT NULL,
    `lineactive` INT NOT NULL,
    `ecbstate` VARCHAR(20) NOT NULL,
    `theme` VARCHAR(20) NOT NULL,
    `tacktime` INT NOT NULL,
    `workcenters` TEXT NOT NULL,
    `status` VARCHAR(20) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Master Finished Goods (unique fgtype+lotinv)
CREATE TABLE IF NOT EXISTS `masterfgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `mattype` VARCHAR(10) NOT NULL,
    `matdesc` VARCHAR(50) NOT NULL,
    `fgtype` VARCHAR(20) NOT NULL,
    `aging_tipes_id` INT NOT NULL DEFAULT 0,
    `kdbar` VARCHAR(20) NOT NULL,
    `warna` VARCHAR(20) NOT NULL,
    `lotinv` VARCHAR(20) NOT NULL,
    `attrib` VARCHAR(100) NOT NULL,
    `category` VARCHAR(20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `masterfgs_fgtype_lotinv_unique` (`fgtype`, `lotinv`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Master Semi-Finished Goods (unique plant/mattype/sfgtype)
CREATE TABLE IF NOT EXISTS `mastersfgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `plant` VARCHAR(10) NOT NULL,
    `mattype` VARCHAR(10) NOT NULL,
    `matdesc` VARCHAR(50) NOT NULL,
    `sfgtype` VARCHAR(20) NOT NULL,
    `sfgdesc` VARCHAR(50) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `mastersfgs_plant_mattype_sfgtype_unique` (`plant`, `mattype`, `sfgtype`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Compressors
CREATE TABLE IF NOT EXISTS `compressors` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ctype` VARCHAR(4) NOT NULL DEFAULT "",
    `merk` VARCHAR(20) NOT NULL DEFAULT "",
    `type` VARCHAR(20) NOT NULL DEFAULT "",
    `itemcode` VARCHAR(20) NOT NULL DEFAULT "",
    `force_scan` INT UNSIGNED NOT NULL DEFAULT 1,
    `familycode` VARCHAR(20) NOT NULL DEFAULT "",
    `status` VARCHAR(20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Compressor reference to FG prefix
CREATE TABLE IF NOT EXISTS `comprefgs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ctype` VARCHAR(4) NOT NULL DEFAULT "",
    `barcode` VARCHAR(20) NOT NULL DEFAULT "",
    `status` VARCHAR(20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Purchase orders (ctype memakai kode kompresor 4 karakter)
CREATE TABLE IF NOT EXISTS `ecbpos` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `workcenter` VARCHAR(20) NOT NULL,
    `po` VARCHAR(20) NOT NULL,
    `sn` VARCHAR(20) NOT NULL,
    `ctype` VARCHAR(4) NOT NULL DEFAULT "",
    `updated_by` INT UNSIGNED NOT NULL,
    `status` VARCHAR(20) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Production log
CREATE TABLE IF NOT EXISTS `ecbdatas` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tgl` DATE NOT NULL DEFAULT "1920-01-01",
    `jam` TIME NOT NULL DEFAULT "00:00:00",
    `wc` VARCHAR(20) NOT NULL DEFAULT "",
    `prdline` VARCHAR(20) NOT NULL DEFAULT "",
    `ctgr` VARCHAR(255) NOT NULL DEFAULT "",
    `sn` VARCHAR(25) NOT NULL DEFAULT "",
    `fgtype` VARCHAR(20) NOT NULL DEFAULT "",
    `spc` VARCHAR(20) NOT NULL DEFAULT "",
    `comptype` VARCHAR(20) NOT NULL DEFAULT "",
    `compcode` VARCHAR(30) NOT NULL DEFAULT "",
    `po` VARCHAR(20) NOT NULL DEFAULT "",
    `status` VARCHAR(20) NOT NULL DEFAULT "",
    `sendsts` VARCHAR(20) NOT NULL DEFAULT "",
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `ecbdatas_1_index` (`tgl`, `jam`, `sn`),
    KEY `ecbdatas_2_index` (`ctgr`, `prdline`),
    KEY `ecbdatas_3_index` (`fgtype`),
    KEY `ecbdatas_4_index` (`comptype`),
    KEY `ecbdatas_5_index` (`sendsts`),
    KEY `ecbdatas_6_index` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Seed contoh minimal
INSERT INTO `masterfgs` (mattype, matdesc, fgtype, aging_tipes_id, kdbar, warna, lotinv, attrib, category, created_at, updated_at) VALUES
  ('FG', 'REF Aurora 180L', 'REF', 1, 'RF11', 'Silver', 'IDN0', 'A', 'REFRIG', NOW(), NOW()),
  ('FG', 'REF Boreal 220L', 'REF', 1, 'RF12', 'Black', 'REF-LOT-B', 'B', 'REFRIG', NOW(), NOW()),
  ('FG', 'REF Cobalt 260L', 'REF', 1, 'RF21', 'Gray', 'REF-LOT-C', 'C', 'REFRIG', NOW(), NOW()),
  ('FG', 'PWM Breeze 9kg', 'PWM', 1, 'PW10', 'White', 'IDN0', 'A', 'WASH', NOW(), NOW()),
  ('FG', 'PWM Cyclone 12kg', 'PWM', 1, 'PW11', 'Graphite', 'PWM-LOT-B', 'B', 'WASH', NOW(), NOW());

INSERT INTO `mastersfgs` (plant, mattype, matdesc, sfgtype, sfgdesc, created_at, updated_at) VALUES
  ('PLT1', 'SG01', 'REF Subassy D', 'SFG-REF-A', 'Door & gasket', NOW(), NOW()),
  ('PLT1', 'SG02', 'REF Subassy C', 'SFG-REF-B', 'Foam & piping', NOW(), NOW()),
  ('PLT2', 'SG03', 'PWM Subassy', 'SFG-PWM-A', 'Tub & pulsator', NOW(), NOW()),
  ('PLT2', 'SG04', 'PWM Harness', 'SFG-PWM-B', 'Harness & PCB', NOW(), NOW());

INSERT INTO `compressors` (ctype, merk, type, itemcode, force_scan, familycode, status, created_at, updated_at) VALUES
  ('RFA1', 'Panasonic', 'XA-10', 'CMPA10', 0, 'REF100', 'ACTIVE', NOW(), NOW()),
  ('RFB1', 'LG', 'BL-22', 'CMPB22', 0, 'REF110', 'ACTIVE', NOW(), NOW()),
  ('RFC1', 'Samsung', 'CS-33', 'CMPC33', 1, 'REF120', 'ACTIVE', NOW(), NOW()),
  ('PWM1', 'Toshiba', 'TD-12', 'CMPD12', 0, 'PWM10', 'ACTIVE', NOW(), NOW()),
  ('PWM2', 'Hitachi', 'HE-50', 'CMPE50', 1, 'PWM11', 'ACTIVE', NOW(), NOW());

INSERT INTO `comprefgs` (ctype, barcode, status, created_at, updated_at) VALUES
  ('RFA1', 'RF11', 'OK', NOW(), NOW()),
  ('RFA1', 'RF21', 'OK', NOW(), NOW()),
  ('RFB1', 'RF12', 'OK', NOW(), NOW()),
  ('RFC1', 'RF21', 'OK', NOW(), NOW()),
  ('PWM1', 'PW10', 'OK', NOW(), NOW()),
  ('PWM2', 'PW11', 'OK', NOW(), NOW());

INSERT INTO `ecbstations` (ipaddress, location, mode, linetype, lineids, lineactive, ecbstate, theme, tacktime, workcenters, status, created_at, updated_at) VALUES
  ('192.168.10.10', 'REF Remote Cluster', 'LIVE', 'refrig-po-double', 'REF A,REF B', 0, 'READY', 'Minimal Night', 55, 'WC-REF-A,WC-REF-B', 'ACTIVE', NOW(), NOW()),
  ('192.168.10.20', 'PWM Remote', 'LIVE', 'sn-only-single', 'PWM 1', 0, 'READY', 'Minimal Night', 60, 'WC-PWM', 'ACTIVE', NOW(), NOW());

INSERT INTO `ecbpos` (workcenter, po, sn, ctype, updated_by, status, created_at, updated_at) VALUES
  ('REF A', 'PO-REMOTE-REF-A', 'RF1100009999', 'RFA1', 9, 'OPEN', NOW(), NOW()),
  ('REF B', 'PO-REMOTE-REF-B', 'RF1200008888', 'RFB1', 9, 'OPEN', NOW(), NOW()),
  ('PWM 1', 'PO-REMOTE-PWM', 'PW1000007777', 'PWM1', 9, 'OPEN', NOW(), NOW()),
  ('REF A', 'PO-REMOTE-REF-A2', 'RF1199880999', 'RFA1', 9, 'OPEN', NOW(), NOW()),
  ('PWM 1', 'PO-REMOTE-PWM2', 'PW1099880999', 'PWM1', 9, 'OPEN', NOW(), NOW());

INSERT INTO `ecbdatas` (tgl, jam, wc, prdline, ctgr, sn, fgtype, spc, comptype, compcode, po, status, sendsts, created_at, updated_at) VALUES
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '06:50:00', 'WC-REF-A', 'REF A', 'REF', 'RF1098001234', 'REF', 'SPCREFREM1', 'RFA1', 'CMPA10-REM1', 'PO-REMOTE-REF-A', 'PASS', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '06:55:00', 'WC-REF-B', 'REF B', 'REF', 'RF1298001235', 'REF', 'SPCREFREM2', 'RFB1', 'CMPB22-REM2', 'PO-REMOTE-REF-B', 'FAIL', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 2 DAY), '07:05:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1098001236', 'PWM', 'SPCPWMREM3', 'PWM1', 'CMPD12-REM3', '', 'PASS', 'SENT', NOW(), NOW());
