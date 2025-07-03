require 'pg'

Before('@database') do
  config = YAML.load_file(TestAra.instance.config_dir.join('database.yml'))["test"]
  config["dbname"] = config.delete("name")
  @connection = PG.connect config
end

After('@database') do
  # Truncate all tables
  @connection.exec(
    "DO $$DECLARE statements CURSOR FOR
      SELECT table_name FROM information_schema.tables
      WHERE table_schema='public' AND table_name NOT IN ('gorp_migrations', 'ar_internal_metadata', 'schema_migrations', 'geometry_columns', 'geography_columns');
    BEGIN
      FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.table_name) || ' RESTART IDENTITY CASCADE;';
      END LOOP;
    END$$;")
end
