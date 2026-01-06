/*
    file:           cmd/migrate/migrations/20251028043402_NavigationSeed.down.sql
    description:    Rollback seed by clearing navigations records.
    created:        220711663@students.uajy.ac.id 04-11-2025
*/
TRUNCATE TABLE `navigations`;
