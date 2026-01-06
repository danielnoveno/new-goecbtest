/*
    file:           cmd/migrate/migrations/20251028062754_create_ecbdatas_table.down.sql
    description:    Rollback script removing the ecbdatas production capture table.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
DROP TABLE IF EXISTS `ecbdatas` ;
