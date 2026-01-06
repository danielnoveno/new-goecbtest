/*
    file:           cmd/migrate/migrations/20251028045833_create_ecbstates_table.down.sql
    description:    Rollback script dropping the ECB states table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `ecbstates`;
