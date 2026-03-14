#!/usr/bin/env bash
set -euo pipefail

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_APP_USER" -d "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE comments (
		id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		user_id bigint NOT NULL,
		content TEXT NOT NULL
	);
EOSQL
