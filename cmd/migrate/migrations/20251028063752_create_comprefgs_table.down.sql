/*
    file:           cmd/migrate/migrations/20251028063752_create_comprefgs_table.down.sql
    description:    Rollback script dropping comprefgs and related ecbdatas indexes.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `comprefgs`;

ALTER TABLE ecbdatas DROP INDEX ecbdatas_1_index;
ALTER TABLE ecbdatas DROP INDEX ecbdatas_2_index;
ALTER TABLE ecbdatas DROP INDEX ecbdatas_3_index;
ALTER TABLE ecbdatas DROP INDEX ecbdatas_4_index;
ALTER TABLE ecbdatas DROP INDEX ecbdatas_5_index;
ALTER TABLE ecbdatas DROP INDEX ecbdatas_6_index;
