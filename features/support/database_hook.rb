require 'pg'

$database = 'edwig_test'

After('@database') do
  # Truncate all tables
  conn = PG.connect dbname: $database, user: ENV["POSTGRESQL_ENV_POSTGRES_USER"], password: ENV["POSTGRESQL_ENV_POSTGRES_PASSWORD"]
  conn.exec(
    "DO $$DECLARE statements CURSOR FOR
      SELECT table_name FROM information_schema.tables
      WHERE table_schema='public' AND table_name NOT IN ('gorp_migrations', 'ar_internal_metadata', 'schema_migrations');
    BEGIN
      FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.table_name) || ' CASCADE;';
      END LOOP;
    END$$;")
end
