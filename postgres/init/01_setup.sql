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
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(), -- noqa: CP03
  title text NOT NULL,
  content text,
  done bool NOT NULL DEFAULT FALSE,
  version int NOT NULL DEFAULT 1,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
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
