begin;

CREATE TABLE IF NOT EXISTS app_settings (
    id SERIAL PRIMARY KEY,
    platform_path TEXT NOT NULL
);

COMMENT ON TABLE app_settings IS 'Настройки приложения';

commit;