begin;

CREATE TYPE db_state AS ENUM ('new', 'analyzing', 'done', 'error');
CREATE TYPE job_state AS ENUM ('new', 'in_progress', 'error', 'done');

CREATE TABLE IF NOT EXISTS database_info (
            id SERIAL PRIMARY KEY,
            last_check TIMESTAMPTZ,
            name TEXT NOT NULL,
            status db_state NOT NULL DEFAULT 'new',
            connection_string TEXT NOT NULL,
            username TEXT,
            password TEXT
);

CREATE TABLE IF NOT EXISTS conf_info (
    id SERIAL PRIMARY KEY,
    database_id integer REFERENCES database_info(id) ON DELETE CASCADE UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    version TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS extensions_info (
    conf_id integer REFERENCES conf_info(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    version TEXT NOT NULL,
    purpose TEXT NOT NULL,
    PRIMARY KEY (conf_id, name)
);

CREATE TABLE IF NOT EXISTS jobs (
        database_id integer REFERENCES database_info(id) ON DELETE CASCADE PRIMARY KEY,
        next_check TIMESTAMPTZ,
        cron TEXT NOT NULL,
        status job_state NOT NULL DEFAULT 'new'
);


COMMENT ON TABLE database_info IS 'Настройки подключения к БД 1с';
COMMENT ON TABLE conf_info IS 'Информация об конфигурациях';
COMMENT ON TABLE extensions_info IS 'Информация об расширениях';
COMMENT ON TABLE jobs IS 'Задания';

commit;