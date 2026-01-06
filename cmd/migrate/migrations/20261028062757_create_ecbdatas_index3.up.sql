/*
    file:           cmd/migrate/migrations/20261028062757_create_ecbdatas_index3.up.sql
    description:    Migration creating fgtype index to speed finished goods lookups.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_3_index ON `ecbdatas` (fgtype);
