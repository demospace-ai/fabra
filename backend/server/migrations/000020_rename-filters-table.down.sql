ALTER TABLE IF EXISTS event_filters RENAME TO step_filters;
ALTER TABLE IF EXISTS step_filters RENAME COLUMN event_id TO step_id;
ALTER INDEX IF EXISTS event_filters_pkey RENAME TO step_filters_pkey;
ALTER INDEX IF EXISTS event_filters_analysis_id_idx RENAME TO step_filters_analysis_id_idx;
ALTER INDEX IF EXISTS event_filters_step_id_idx RENAME TO step_filters_step_id_idx;

ALTER TABLE IF EXISTS events RENAME TO funnel_steps;
ALTER TABLE IF EXISTS funnel_steps RENAME COLUMN name TO step_name;
ALTER INDEX IF EXISTS events_pkey RENAME TO funnel_steps_pkey;
ALTER INDEX IF EXISTS events_analysis_id_idx RENAME TO funnel_steps_analysis_id_idx;