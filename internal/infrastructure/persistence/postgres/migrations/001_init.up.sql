CREATE TABLE IF NOT EXISTS teams (
                       name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
                       id TEXT PRIMARY KEY,
                       username TEXT NOT NULL,
                       team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE RESTRICT,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pull_requests (
                               id TEXT PRIMARY KEY,
                               name TEXT NOT NULL,
                               author_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                               status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
                               merged_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
                              pr_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
                              user_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                              PRIMARY KEY(pr_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id);