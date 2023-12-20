CREATE TABLE IF NOT EXISTS assignment_rating_cache (
    assignment_id int PRIMARY KEY,
    avg_rating float,
    expires_at timestamptz not null
);

create index on assignment_rating_cache (assignment_id);
create index on assignment_rating_cache (expires_at);
