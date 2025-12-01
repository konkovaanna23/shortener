
CREATE TABLE IF NOT EXISTS urls.users (
	id varchar(200),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS urls.links_users_xmap (
	user_id varchar(200),
	url_id int
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_uq_links_users_xmap ON urls.links_users_xmap (url_id);