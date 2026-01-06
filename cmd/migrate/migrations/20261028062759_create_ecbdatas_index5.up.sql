/*
    file:           cmd/migrate/migrations/20261028062759_create_ecbdatas_index5.up.sql
    description:    Migration creating send status index for ecbdatas records.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_5_index ON `ecbdatas` (sendsts);
