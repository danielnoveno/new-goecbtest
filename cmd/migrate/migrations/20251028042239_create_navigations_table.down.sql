/*
    file:           cmd/migrate/migrations/20251028042239_create_navigations_table.down.sql
    description:    Rollback script dropping the navigations menu table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `navigations`;
