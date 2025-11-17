DO $$
BEGIN

    IF NOT EXISTS(
        SELECT schema_name
          FROM information_schema.schemata
          WHERE schema_name = 'urls'
      )
    THEN
      EXECUTE 'CREATE SCHEMA urls';
    END IF;

END
$$;

CREATE TABLE IF NOT EXISTS urls.links (
	uuid int not null,
	short_url varchar(255) not null,
	original_url text not null,
    created_at TIMESTAMP DEFAULT NOW(),
	PRIMARY KEY (short_url)
)
