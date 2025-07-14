CREATE TABLE IF NOT EXISTS diaries (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    date DATE NOT NULL,
    mental INTEGER NOT NULL,
    diary TEXT NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);