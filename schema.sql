CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
           NEW.updated_at = now();
           RETURN NEW;
        END;
        $$ language 'plpgsql';



CREATE TABLE users (
   id serial primary key,
   github_access_token varchar(255),
   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE
        ON users FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();


CREATE TYPE hosting_provider AS ENUM ('github');

CREATE TABLE issues (
   id serial primary key,
   original_url varchar(255),

   hoster hosting_provider,
   project varchar(255),
   repo varchar(255),
   identifier varchar(20),

   coinbase_button_code varchar(100),

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
