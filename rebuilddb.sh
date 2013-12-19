dropdb bountyforcode
createdb bountyforcode
psql bountyforcode < schema.sql

dropdb bountyforcode_test
createdb bountyforcode_test
psql bountyforcode_test < schema.sql
psql bountyforcode_test < test.sql
