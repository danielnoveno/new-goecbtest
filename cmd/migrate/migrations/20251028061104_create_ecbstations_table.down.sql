/*
    file:           cmd/migrate/migrations/20251028061104_create_ecbstations_table.down.sql
    description:    Rollback script dropping the ECB stations table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `ecbstations` ;
