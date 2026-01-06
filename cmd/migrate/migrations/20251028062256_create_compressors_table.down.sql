/*
    file:           cmd/migrate/migrations/20251028062256_create_compressors_table.down.sql
    description:    Rollback script dropping the compressors reference table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `compressors` ;
