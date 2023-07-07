ALTER TABLE IF EXISTS step_filters RENAME TO event_filters;
ALTER TABLE IF EXISTS event_filters RENAME COLUMN step_id TO event_id;
ALTER INDEX IF EXISTS step_filters_pkey RENAME TO event_filters_pkey;
ALTER INDEX IF EXISTS step_filters_analysis_id_idx RENAME TO event_filters_analysis_id_idx;
ALTER INDEX IF EXISTS step_filters_step_id_idx RENAME TO event_filters_step_id_idx;

ALTER TABLE IF EXISTS funnel_steps RENAME TO events;
ALTER TABLE IF EXISTS events RENAME COLUMN step_name TO name;
ALTER INDEX IF EXISTS funnel_steps_pkey RENAME TO events_pkey;
ALTER INDEX IF EXISTS funnel_steps_analysis_id_idx RENAME TO events_analysis_id_idx;