begin;

CREATE TABLE IF NOT EXISTS metadata (
            conf_id integer REFERENCES conf_info(id) ON DELETE CASCADE PRIMARY KEY,
            struct JSONB NOT NULL
);

COMMENT ON TABLE metadata IS 'Хранение структуры метаданных расширений по конфигурации';

commit;