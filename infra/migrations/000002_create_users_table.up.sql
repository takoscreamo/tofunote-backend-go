CREATE TABLE users (
  id UUID PRIMARY KEY,
  nickname VARCHAR(255),
  provider VARCHAR(50),
  provider_id VARCHAR(255),
  is_guest BOOLEAN DEFAULT TRUE,
  refresh_token VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 