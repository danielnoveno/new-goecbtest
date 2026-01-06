/*
    file:           cmd/migrate/migrations/20261028062756_create_ecbdatas_index2.up.sql
    description:    Migration adding category and production line index for ecbdatas queries.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_2_index ON `ecbdatas` (ctgr, prdline);
