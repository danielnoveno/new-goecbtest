/*
    file:           cmd/migrate/migrations/20251028064016_create_ecbpos_table.down.sql
    description:    Rollback script dropping the ECB PO table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `ecbpos`;
