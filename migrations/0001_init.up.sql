CREATE TABLE IF NOT EXISTS grant_projects (id text primary key, code text unique not null, name text not null, year integer not null, annual_cap_cents bigint not null, version integer not null default 1);
CREATE TABLE IF NOT EXISTS expense_claims (id text primary key, project_id text not null, beneficiary_id text not null, amount_cents bigint not null, status text not null, idempotency_key text unique not null);
CREATE TABLE IF NOT EXISTS settlements (id text primary key, claim_id text not null, subsidy_cents bigint not null, status text not null, version integer not null default 1);
