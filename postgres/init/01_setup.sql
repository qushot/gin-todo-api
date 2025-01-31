CREATE FUNCTION set_version_updated_at() RETURNS TRIGGER AS -- noqa: CP03
$$
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        NEW.version := OLD.version + 1;
        NEW.updated_at := NOW();
        return NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS todo (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid() -- noqa: CP03
  , title TEXT NOT NULL
  , content TEXT
  , done BOOL NOT NULL DEFAULT FALSE
  , version INT NOT NULL DEFAULT 1
  , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
COMMENT ON TABLE todo IS 'ToDo';
COMMENT ON COLUMN todo.id IS 'ID';
COMMENT ON COLUMN todo.title IS 'タイトル';
COMMENT ON COLUMN todo.content IS '内容';
COMMENT ON COLUMN todo.done IS '完了フラグ';
COMMENT ON COLUMN todo.version IS 'バージョン';
COMMENT ON COLUMN todo.created_at IS '作成日時';
COMMENT ON COLUMN todo.updated_at IS '更新日時';

CREATE OR REPLACE TRIGGER trg_todo_version_updated_at
BEFORE UPDATE ON todo FOR EACH ROW EXECUTE PROCEDURE set_version_updated_at(); -- noqa: CP03
COMMENT ON TRIGGER trg_todo_version_updated_at ON todo IS 'バージョンと更新日時を更新するトリガー';


CREATE TABLE IF NOT EXISTS sample (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid() -- noqa: CP03
  , todo_id UUID REFERENCES todo (id)
  , version INT NOT NULL DEFAULT 1
  , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
COMMENT ON TABLE sample IS 'sample';
COMMENT ON COLUMN sample.id IS 'ID';
COMMENT ON COLUMN sample.todo_id IS 'ToDo ID';
COMMENT ON COLUMN sample.version IS 'バージョン';
COMMENT ON COLUMN sample.created_at IS '作成日時';
COMMENT ON COLUMN sample.updated_at IS '更新日時';

CREATE OR REPLACE TRIGGER trg_sample_version_updated_at
BEFORE UPDATE ON sample FOR EACH ROW EXECUTE PROCEDURE set_version_updated_at(); -- noqa: CP03
COMMENT ON TRIGGER trg_sample_version_updated_at ON sample IS 'バージョンと更新日時を更新するトリガー';
