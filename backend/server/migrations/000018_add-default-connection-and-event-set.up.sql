ALTER TABLE organizations ADD default_data_connection_id BIGINT REFERENCES data_connections(id);
ALTER TABLE organizations ADD default_event_set_id BIGINT REFERENCES event_sets(id);