--Tables
CREATE TABLE IF NOT EXISTS pasty (pasta_id SERIAL PRIMARY KEY, pasta_name text, pasta_body text);
CREATE TABLE IF NOT EXISTS pasty_tag_relation (pasta_id integer, tag_id integer);
CREATE TABLE IF NOT EXISTS tags (tag_id SERIAL PRIMARY KEY, tag_root text, tag_name text);
CREATE TABLE IF NOT EXISTS users (credentials text, hashed_password text);

--Functions

CREATE OR REPLACE FUNCTION INSERT_PASTA (pasta_name text, 
								pasta_body text,
								tag_names text ARRAY)
RETURNS INT AS
$$
DECLARE
tag text;
tagid int;
pastaid int;
BEGIN
	INSERT INTO pasty (pasta_name, pasta_body) VALUES (pasta_name, pasta_body) 
		RETURNING pasta_id INTO pastaid;
	FOREACH tag IN ARRAY tag_names
	LOOP
	BEGIN
		SELECT tag_id INTO tagid FROM tags WHERE tag_name = tag;
		INSERT INTO pasty_tag_relation VALUES (pastaid, tagid);
	EXCEPTION
        WHEN NO_DATA_FOUND THEN
			CONTINUE;
	END;
	END LOOP;
	RETURN pastaid;
END
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION GET_RECORDS () 
	RETURNS TABLE(pastaid int, 
		pastaname text,
		tags text ARRAY)
	AS
$$
BEGIN
	RETURN QUERY 
		SELECT pasta_id, pasta_name, 
		ARRAY(
			SELECT tag_name FROM tags INNER JOIN pasty_tag_relation
			ON pasty_tag_relation.tag_id=tags.tag_id 
				WHERE pasty_tag_relation.pasta_id=pasty.pasta_id)
		FROM pasty;
END
$$
LANGUAGE plpgsql;