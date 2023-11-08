CREATE TABLE IF NOT EXISTS movies (id VARCHAR(255),title VARCHAR(255),description TEXT,director VARCHAR(255));

CREATE TABLE IF NOT EXISTS ratings (record_id VARCHAR(255),rating_type VARCHAR(255),user_id VARCHAR(255),value INT);

-- docker exec -i movieexample_db1 mysql movieexample -h localhost -P 3306 --protocol=tcp -uroot -ppassword < schema/schema.sql

--docker exec -i movieexample_db1 mysql movieexample -h localhost -P 3306 --protocol=tcp -uroot -ppassword -e "SHOW tables"