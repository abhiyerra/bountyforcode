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


CREATE TABLE issues (
   id serial primary key,
   original_url varchar(255) default '',

   hoster varchar(20) not null default '',
   project varchar(255) not null default '',
   repo varchar(255) not null default '',
   identifier varchar(20) not null default '',

   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_issues_updated_at BEFORE UPDATE
        ON issues FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();

CREATE TABLE bounties (
   id serial primary key,
   user_id serial references users(id),
   issue_id serial references issues(id),
   amount float not null default 0.0,

   coinbase_button_code varchar(255) not null default '',
   coinbase_order_id varchar(100) not null default '',
   coinbase_total_btc int not null default 0,
   coinbase_currency_iso varchar(5) default '',

   status varchar(10) not null default 'new',

   created_at timestamp default now(),
   updated_at timestamp default now()
);

CREATE TRIGGER update_bounties_updated_at BEFORE UPDATE
        ON issues FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at_column();
