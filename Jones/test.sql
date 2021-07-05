CREATE TABLE test_0702 (id INT, NAME VARCHAR(200),age VARCHAR(200),state VARCHAR(200) );

ALTER TABLE test_0702 ADD COLUMN (state VARCHAR(200))



CREATE INDEX id_name_age ON test_0702(id,NAME,age);

CREATE INDEX id_name ON test_0702(id,NAME);

CREATE INDEX id_index ON test_0702(id);



INSERT INTO test_0702 VALUE(1,"zhanbgsan",21,"中国");
INSERT INTO test_0702 VALUE(2,"zhanbgsan1",22,"中国");
INSERT INTO test_0702 VALUE(3,"lisi",22,"中国");
INSERT INTO test_0702 VALUE(4,"wangwu",21,"中国");
INSERT INTO test_0702 VALUE(5,"zhaoliu",21,"中国");
INSERT INTO test_0702 VALUE(6,"caiqi",21,"中国");


EXPLAIN SELECT * FROM test_0702 WHERE id = 1 AND age="23"

DROP TABLE test_0702


SHOW INDEX  FROM test_0702
