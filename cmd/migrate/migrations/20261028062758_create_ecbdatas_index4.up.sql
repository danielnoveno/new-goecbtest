/*
    file:           cmd/migrate/migrations/20261028062758_create_ecbdatas_index4.up.sql
    description:    Migration adding component type index for ecbdatas.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_4_index ON `ecbdatas` (comptype);
