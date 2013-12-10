CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
           NEW.updated_at = now();
           RETURN NEW;
        END;
        $$ language 'plpgsql';


CREATE TYPE hosting_provider AS ENUM ('github');

CREATE TABLE users (
   id serial primary key,
   username varchar(255),
   password varchar(255),
   salt varchar(255),
   github_identifier varchar(255),
   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE
        ON users FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();


CREATE TABLE issues (
   id serial primary key,
   original_url varchar(255),
   hoster hosting_provider,
   repo varchar(255),
   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_issues_updated_at BEFORE UPDATE
        ON issues FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();

CREATE TYPE bounty_state as ENUM ('open', 'paid', 'closed', 'cancelled');

CREATE TABLE bounties (
   id serial primary key,
   user_id serial references users(id),
   issue_id serial references issues(id),
   amount float,
   transaction_status bounty_state,
   expires_at timestamp,
   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_bounties_updated_at BEFORE UPDATE
        ON issues FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();
