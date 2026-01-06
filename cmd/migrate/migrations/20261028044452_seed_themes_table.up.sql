/*
    file:           cmd/migrate/migrations/20261028044452_seed_themes_table.up.sql
    description:    Seed the new Fyne-friendly theme palette store with initial ensembles.
    created:        codex-agent
*/
INSERT INTO `themes` (
    nama,
    keterangan,
    color_background,
    color_foreground,
    color_text,
    color_button,
    color_disabled,
    color_error,
    color_focus,
    color_hover,
    color_input_background,
    color_placeholder,
    color_primary,
    color_scrollbar,
    color_selection,
    color_navbar,
    color_footer,
    header_start,
    header_end,
    accent,
    created_at,
    updated_at
) VALUES
('Minimal Night', 'Palet gelap minimal dengan aksen biru muda.', '#101214', '#E4E7EC', '#E4E7EC', '#3F51B5', '#5F6670', '#E57373', '#7986CB', '#2C2F33', '#1E2228', '#8F8F99', '#3F51B5', '#2C2F33', '#1E88E5', '#121318', '#121318', '#121318', '#121318', '#4DD0E1', NOW(), NOW()),
('Minimal Light', 'Palet terang minimal menonjolkan kontras biru.', '#F5F6F7', '#1F2733', '#1F2733', '#1976D2', '#B0BEC5', '#E53935', '#90CAF9', '#E3F2FD', '#FFFFFF', '#9E9E9E', '#1976D2', '#CFD8DC', '#BBDEFB', '#ECEFF1', '#ECEFF1', '#ECEFF1', '#ECEFF1', '#1976D2', NOW(), NOW()),
('Minimal Mono', 'Nuansa abu-abu dengan sentuhan netral untuk fokus.', '#1A1A1A', '#F5F5F5', '#F5F5F5', '#757575', '#424242', '#EF5350', '#BDBDBD', '#2E2E2E', '#0F0F0F', '#9E9E9E', '#757575', '#111111', '#424242', '#121212', '#121212', '#121212', '#121212', '#90A4AE', NOW(), NOW()),
('Minimal Warm', 'Kombinasi hangat tembaga dan kayu untuk sentuhan humanis.', '#20120B', '#FFE5CB', '#FFE5CB', '#FF7043', '#8D6E63', '#FF5252', '#FFAB91', '#3E2723', '#2C160C', '#BCAAA4', '#FF7043', '#311711', '#FF8A65', '#2E1A15', '#2E1A15', '#2E1A15', '#2E1A15', '#FFAB40', NOW(), NOW()),
('Minimal Forest', 'Palet hijau tua dengan aksen lembut untuk suasana tenang.', '#0B1B12', '#C5E1A5', '#C5E1A5', '#66BB6A', '#356B34', '#E53935', '#A5D6A7', '#1A361F', '#0F140D', '#8E9F8A', '#66BB6A', '#12230E', '#81C784', '#112014', '#112014', '#112014', '#112014', '#4CAF50', NOW(), NOW());
