-- +migrate Up notransaction
CREATE INDEX CONCURRENTLY ON app_version_events(app_version_id);

-- +migrate Down
DROP INDEX app_version_events_app_version_id_idx;