/*
    file:           cmd/migrate/migrations/20251028062102_create_mastersfgs_table.down.sql
    description:    Rollback script dropping the semi-finished goods table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `mastersfgs`;
