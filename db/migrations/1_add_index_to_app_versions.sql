-- +migrate Up notransaction
CREATE INDEX CONCURRENTLY ON app_versions(app_id);

-- +migrate Down
DROP INDEX app_versions_app_id_idx;