-- +migrate Up notransaction
CREATE INDEX CONCURRENTLY ON app_settings(app_id);

-- +migrate Down
DROP INDEX app_settings_app_id_idx;