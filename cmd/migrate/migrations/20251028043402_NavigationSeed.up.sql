/*
    file:           cmd/migrate/migrations/20251028043402_NavigationSeed.up.sql
    description:    Seed data inserting default top-level navigation entries.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
INSERT INTO `navigations` (parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at) VALUES
 (NULL, 'bolt', 'ECB Test', 'Panel pengujian ECB terpadu', '/ecb-test', 'ecb-test', 1, 1, NOW(), NOW()),
 (NULL, 'settings', 'Setting', 'Pengaturan sistem simulasi', '/settings', 'settings', 1, 2, NOW(), NOW()),
 (NULL, 'performance', 'Maintenance', 'Pantau status dan pemeliharaan', '/maintenance', 'maintenance', 1, 3, NOW(), NOW()),
 (NULL, 'info', 'About', 'Tentang aplikasi dan tim pengembang', '/about', 'about', 1, 4, NOW(), NOW()),
 (NULL, 'power-off', 'Shutdown', 'Kirim perintah shutdown perangkat', '/shutdown', 'shutdown', 1, 5, NOW(), NOW()),
 (NULL, 'reboot', 'Reboot', 'Reset sesi dan subsystem', '/reboot', 'reboot', 1, 6, NOW(), NOW());
