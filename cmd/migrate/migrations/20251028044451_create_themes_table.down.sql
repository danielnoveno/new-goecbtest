/*
    file:           cmd/migrate/migrations/20251028044451_create_themes_table.down.sql
    description:    Rollback script removing the themes lookup table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `themes`;
