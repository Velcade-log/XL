-- ============================================
-- 插入初始用户账号
-- 执行前请确保已经执行过 init.sql 建表
-- ============================================

-- 账号1：用户名 15813272880，密码 15813272880xl
-- 账号2：用户名 13721999775，密码 13721999775xs
INSERT INTO users (username, password_hash) VALUES
    ('15813272880', '$2a$10$ZH3jRfKJjlFD9XqJTEdCvevDshkw3bOUcvlNnXWirA4QEtTDaMAKW'),
    ('13721999775', '$2a$10$NvtgvXJCvLaEEm1Gqo5unuLaA1ukdGj80KZoa7fnbQm9SDIGQiwde');
