/*
    file:           cmd/migrate/migrations/20261028062755_create_ecbdatas_index1.up.sql
    description:    Migration adding composite index on ecbdatas timestamp and serial.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
CREATE INDEX ecbdatas_1_index ON `ecbdatas` (tgl, jam, sn);
