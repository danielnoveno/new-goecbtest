/*
    file:           cmd/migrate/migrations/20261028044451_create_themes_table.up.sql
    description:    Migration that creates the themes table with color tokens for the Fyne UI.
    created:        codex-agent
*/
CREATE TABLE IF NOT EXISTS `themes` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(60) NOT NULL,
    `keterangan` TEXT NOT NULL,
    `color_background` VARCHAR(7) NOT NULL,
    `color_foreground` VARCHAR(7) NOT NULL,
    `color_text` VARCHAR(7) NOT NULL,
    `color_button` VARCHAR(7) NOT NULL,
    `color_disabled` VARCHAR(7) NOT NULL,
    `color_error` VARCHAR(7) NOT NULL,
    `color_focus` VARCHAR(7) NOT NULL,
    `color_hover` VARCHAR(7) NOT NULL,
    `color_input_background` VARCHAR(7) NOT NULL,
    `color_placeholder` VARCHAR(7) NOT NULL,
    `color_primary` VARCHAR(7) NOT NULL,
    `color_scrollbar` VARCHAR(7) NOT NULL,
    `color_selection` VARCHAR(7) NOT NULL,
    `color_navbar` VARCHAR(7) NOT NULL,
    `color_footer` VARCHAR(7) NOT NULL,
    `header_start` VARCHAR(7) NOT NULL,
    `header_end` VARCHAR(7) NOT NULL,
    `accent` VARCHAR(7) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
