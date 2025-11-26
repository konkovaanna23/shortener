DROP TABLE IF EXISTS urls.links;

DO $$
BEGIN

    IF NOT EXISTS(
        SELECT schema_name
          FROM information_schema.schemata
          WHERE schema_name = 'urls'
      )
    THEN
      DROP SCHEMA urls;
    END IF;

END
$$;
