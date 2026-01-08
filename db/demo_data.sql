-- Dummy data untuk DB lokal `goecbtest` setelah migrate up
-- Jalankan sekali untuk menyiapkan login, master data, PO, dan contoh log ECB.
-- Kolom `ecbpos.ctype` diselaraskan ke kode kompresor 4 karakter agar flow PO tidak gagal.
ALTER TABLE ecbpos MODIFY COLUMN ctype VARCHAR(4) NOT NULL;

INSERT INTO ecbconfigs (section, variable, value, ordering, created_at, updated_at) VALUES
  ('general', 'mode', 'simulateAll', '001', NOW(), NOW()),
  ('general', 'default_line', 'REF A', '002', NOW(), NOW()),
  ('general', 'line_select_gpio', '23', '003', NOW(), NOW()),
  ('general', 'scan_timeout_ms', '750', '004', NOW(), NOW()),
  ('general', 'line_type', 'refrig-po-double', '005', NOW(), NOW()),
  ('display', 'theme', 'Dark', '006', NOW(), NOW());

INSERT INTO ecbstations (
  ipaddress, location, mode, linetype, lineids, lineactive,
  ecbstate, theme, tacktime, workcenters, status,
  created_at, updated_at
) VALUES
  ('192.168.0.101', 'Line REF Double', 'simulateAll', 'refrig-po-double', 'REF A,REF B', 0, 'READY', 'Dark', 55, 'WC-REF-A,WC-REF-B', 'ACTIVE', NOW(), NOW()),
  ('192.168.0.102', 'Line PWM', 'simulateAll', 'sn-only-single', 'PWM 1', 0, 'READY', 'Dark', 60, 'WC-PWM', 'ACTIVE', NOW(), NOW()),
  ('192.168.0.103', 'Line REF Single', 'simulateAll', 'sn-only-single', 'REF C', 0, 'READY', 'Dark', 50, 'WC-REF-C', 'ACTIVE', NOW(), NOW());

INSERT INTO ecbstates (
  tgl, readstate, ecbstate, created_at, updated_at
) VALUES
  (CURDATE(), 'READY', 'IDLE', NOW(), NOW()),
  (CURDATE(), 'RUNNING', 'RUNNING', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), 'MAINTENANCE', 'PAUSED', NOW(), NOW());

INSERT INTO masterfgs (
  mattype, matdesc, fgtype, aging_tipes_id, kdbar, warna,
  lotinv, attrib, category, created_at, updated_at
) VALUES
  ('FG', 'REF Aurora 180L', 'REF', 1, 'RF11', 'Silver', 'IDN0', 'A', 'REFRIG', NOW(), NOW()),
  ('FG', 'REF Boreal 220L', 'REF', 1, 'RF12', 'Black', 'REF-LOT-B', 'B', 'REFRIG', NOW(), NOW()),
  ('FG', 'REF Cobalt 260L', 'REF', 1, 'RF21', 'Gray', 'REF-LOT-C', 'C', 'REFRIG', NOW(), NOW()),
  ('FG', 'PWM Breeze 9kg', 'PWM', 1, 'PW10', 'White', 'IDN0', 'A', 'WASH', NOW(), NOW()),
  ('FG', 'PWM Cyclone 12kg', 'PWM', 1, 'PW11', 'Graphite', 'PWM-LOT-B', 'B', 'WASH', NOW(), NOW());

INSERT INTO mastersfgs (
  plant, mattype, matdesc, sfgtype, sfgdesc, created_at, updated_at
) VALUES
  ('PLT1', 'SG01', 'REF Subassy D', 'SFG-REF-A', 'Door & gasket', NOW(), NOW()),
  ('PLT1', 'SG02', 'REF Subassy C', 'SFG-REF-B', 'Foam & piping', NOW(), NOW()),
  ('PLT2', 'SG03', 'PWM Subassy', 'SFG-PWM-A', 'Tub & pulsator', NOW(), NOW()),
  ('PLT2', 'SG04', 'PWM Harness', 'SFG-PWM-B', 'Harness & PCB', NOW(), NOW());

INSERT INTO compressors (
  ctype, merk, type, itemcode, force_scan, familycode,
  status, created_at, updated_at
) VALUES
  ('RFA1', 'Panasonic', 'XA-10', 'CMPA10', 0, 'REF100', 'ACTIVE', NOW(), NOW()),
  ('RFB1', 'LG', 'BL-22', 'CMPB22', 0, 'REF110', 'ACTIVE', NOW(), NOW()),
  ('RFC1', 'Samsung', 'CS-33', 'CMPC33', 1, 'REF120', 'ACTIVE', NOW(), NOW()),
  ('PWM1', 'Toshiba', 'TD-12', 'CMPD12', 0, 'PWM10', 'ACTIVE', NOW(), NOW()),
  ('PWM2', 'Hitachi', 'HE-50', 'CMPE50', 1, 'PWM11', 'ACTIVE', NOW(), NOW());

INSERT INTO comprefgs (
  ctype, barcode, status, created_at, updated_at
) VALUES
  ('RFA1', 'RF11', 'OK', NOW(), NOW()),
  ('RFA1', 'RF21', 'OK', NOW(), NOW()),
  ('RFB1', 'RF12', 'OK', NOW(), NOW()),
  ('RFC1', 'RF21', 'OK', NOW(), NOW()),
  ('PWM1', 'PW10', 'OK', NOW(), NOW()),
  ('PWM2', 'PW11', 'OK', NOW(), NOW());

INSERT INTO ecbpos (
  workcenter, po, sn, ctype, updated_by, status,
  created_at, updated_at
) VALUES
  ('REF A', 'PO-REF-A-001', 'RF1100001111', 'RFA1', 1, 'OPEN', NOW(), NOW()),
  ('REF B', 'PO-REF-B-001', 'RF1200002222', 'RFB1', 1, 'OPEN', NOW(), NOW()),
  ('REF B', 'PO-REF-B-002', 'RF2100003333', 'RFC1', 2, 'OPEN', NOW(), NOW()),
  ('REF C', 'PO-REF-C-001', 'RF1199000100', 'RFA1', 2, 'OPEN', NOW(), NOW()),
  ('PWM 1', 'PO-PWM-001', 'PW1000004444', 'PWM1', 2, 'OPEN', NOW(), NOW()),
  ('PWM 1', 'PO-PWM-002', 'PW1100005555', 'PWM2', 2, 'OPEN', NOW(), NOW()),
  ('REF A', 'PO-REF-A-010', 'RF1199880100', 'RFA1', 1, 'OPEN', NOW(), NOW()),
  ('REF B', 'PO-REF-B-010', 'RF1299880100', 'RFB1', 2, 'OPEN', NOW(), NOW()),
  ('PWM 1', 'PO-PWM-003', 'PW1099880100', 'PWM1', 2, 'OPEN', NOW(), NOW());

INSERT INTO ecbdatas (
  tgl, jam, wc, prdline, ctgr, sn, fgtype, spc,
  comptype, compcode, po, status, sendsts, created_at, updated_at
) VALUES
  (CURDATE(), '08:00:00', 'WC-REF-A', 'REF A', 'REF', 'RF1199000001', 'REF', '', '', '', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:05:00', 'WC-REF-A', 'REF A', 'REF', 'RF1199000002', 'REF', 'SPCREF90001', 'RFA1', 'CMPA10-90001', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:10:00', 'WC-REF-A', 'REF A', 'REF', 'RF1199000003', 'REF', 'SPCREF90002', 'RFA1', 'CMPA10-90002', '', 'FAIL', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:15:00', 'WC-REF-A', 'REF A', 'REF', 'RF1199000004', 'REF', 'SPCREF90003', 'RFA1', 'CMPA10-90003', '', 'RETEST', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:20:00', 'WC-REF-A', 'REF A', 'REF', 'RF1100001111', 'REF', 'SPCREFPO001', 'RFA1', 'CMPA10-PO1', 'PO-REF-A-001', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:25:00', 'WC-REF-B', 'REF B', 'REF', 'RF1299000001', 'REF', '', '', '', '', 'PASS', 'SENT', NOW(), NOW()),
  (CURDATE(), '08:30:00', 'WC-REF-B', 'REF B', 'REF', 'RF1200002222', 'REF', 'SPCREFPO002', 'RFB1', 'CMPB22-PO2', 'PO-REF-B-001', 'FAIL', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:35:00', 'WC-REF-B', 'REF B', 'REF', 'RF2100003333', 'REF', 'SPCREFPO003', 'RFC1', 'CMPC33-PO3', 'PO-REF-B-002', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:40:00', 'WC-REF-B', 'REF B', 'REF', 'RF2199000002', 'REF', 'SPCREF91002', 'RFC1', 'CMPC33-91002', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '08:45:00', 'WC-REF-C', 'REF C', 'REF', 'RF1199000101', 'REF', 'SPCREF91003', 'RFA1', 'CMPA10-91003', 'PO-REF-C-001', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '09:00:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1099000001', 'PWM', '', '', '', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '09:05:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1099000002', 'PWM', 'SPCPWM90001', 'PWM1', 'CMPD12-90001', '', 'FAIL', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '09:10:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1000004444', 'PWM', 'SPCPWMPO01', 'PWM1', 'CMPD12-PO1', 'PO-PWM-001', 'PASS', 'SENT', NOW(), NOW()),
  (CURDATE(), '09:15:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1100005555', 'PWM', 'SPCPWMPO02', 'PWM2', 'CMPE50-PO2', 'PO-PWM-002', 'RETEST', 'QUEUED', NOW(), NOW()),
  (CURDATE(), '09:20:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1199000003', 'PWM', 'SPCPWM90002', 'PWM2', 'CMPE50-90002', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '14:00:00', 'WC-REF-A', 'REF A', 'REF', 'RF1198000001', 'REF', 'SPCREFOLD1', 'RFA1', 'CMPA10-OLD', '', 'PASS', 'SENT', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '14:05:00', 'WC-REF-B', 'REF B', 'REF', 'RF1198000002', 'REF', 'SPCREFOLD2', 'RFB1', 'CMPB22-OLD', '', 'FAIL', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '14:10:00', 'WC-REF-B', 'REF B', 'REF', 'RF2198000003', 'REF', 'SPCREFOLD3', 'RFC1', 'CMPC33-OLD', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 1 DAY), '14:15:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1098000001', 'PWM', 'SPCPWMOLD1', 'PWM1', 'CMPD12-OLD', '', 'PASS', 'SENT', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 2 DAY), '07:55:00', 'WC-REF-A', 'REF A', 'REF', 'RF1197000001', 'REF', 'SPCREFLEG1', 'RFA1', 'CMPA10-LEG', '', 'PASS', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 2 DAY), '08:05:00', 'WC-PWM', 'PWM 1', 'PWM', 'PW1197000002', 'PWM', 'SPCPWMLEG2', 'PWM2', 'CMPE50-LEG', '', 'FAIL', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 2 DAY), '08:10:00', 'WC-REF-C', 'REF C', 'REF', 'RF1197000003', 'REF', 'SPCREFLEG3', 'RFA1', 'CMPA10-LEG3', '', 'RETEST', 'QUEUED', NOW(), NOW()),
  (DATE_SUB(CURDATE(), INTERVAL 2 DAY), '08:15:00', 'WC-REF-C', 'REF C', 'REF', 'RF2197000004', 'REF', 'SPCREFLEG4', 'RFC1', 'CMPC33-LEG', '', 'PASS', 'QUEUED', NOW(), NOW());
