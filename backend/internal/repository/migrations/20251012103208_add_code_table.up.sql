begin;

CREATE TABLE IF NOT EXISTS code (
    ext_id integer NOT NULL,
    key text NOT NULL,
    code text NOT NULL,
    PRIMARY KEY(ext_id,key)
);

COMMENT ON TABLE code IS 'Хранение кода переопределенных процедур/функций';

commit;