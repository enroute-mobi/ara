require 'pg'

$database = 'edwig_test'

After('@database') do
  # Truncate all tables
  conn = PG.connect dbname: $database
  conn.exec(
    "DO $$DECLARE statements CURSOR FOR
      SELECT table_name FROM information_schema.tables
      WHERE table_schema='public';
    BEGIN
      FOR stmt IN statements LOOP
        IF stmt.table_name <> 'gorp_migrations' THEN
          EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.table_name) || ' CASCADE;';
        END IF;
      END LOOP;
    END$$;")
end