INSERT INTO referentials (referential_id, slug) VALUES
  ('6ba7b814-9dad-11d1-0000-00c04fd430c8', 'ratp'),
  ('6ba7b814-9dad-11d1-0001-00c04fd430c8', 'keolis')
  ON CONFLICT (referential_id) DO NOTHING;
