/*
    file:           cmd/migrate/migrations/20261028062800_create_ecbdatas_index6.up.sql
    description:    Migration adding status index for ecbdatas to speed filtering.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_6_index ON `ecbdatas` (status);
