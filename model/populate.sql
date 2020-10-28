INSERT INTO referentials (referential_id, slug) VALUES
  ('6ba7b814-9dad-11d1-0000-00c04fd430c8', 'ratp'),
  ('6ba7b814-9dad-11d1-0001-00c04fd430c8', 'keolis')
  ON CONFLICT (referential_id) DO NOTHING;

INSERT INTO partners (id, referential_id, slug, settings, connector_types) VALUES (
  '285ae5cd-dec6-4d84-88a6-0152f8284c4c',
  '6ba7b814-9dad-11d1-0000-00c04fd430c8',
  'ratp',
  '{"remote_url": "http://localhost", "remote_objectid_kind": "Reflex", "remote_credential": "ara_cred"}',
  '["siri-stop-monitoring-request-collector", "siri-check-status-client"]'
) ON CONFLICT (id) DO NOTHING;
