DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_pass;

\ir schema.txt

WITH inserted_users AS (
  INSERT INTO users (username, email, is_active)
  VALUES 
    ('tor', 'tor@hjkl.no', true),
    ('extra', 'extra@hjkl.no', true),
    ('ample', 'ample@hjkl.no', true),
    ('thk', 'thk@hjkl.no', true)
  RETURNING id, username
)
INSERT INTO user_pass (id, password_hash)
VALUES
  ((SELECT id FROM inserted_users WHERE username = 'tor'), 'passw0rd'),
  ((SELECT id FROM inserted_users WHERE username = 'extra'), 'plain_pass'),
  ((SELECT id FROM inserted_users WHERE username = 'ample'), 'dont_do_it'),
  ((SELECT id FROM inserted_users WHERE username = 'thk'), 'backdoor');

-- INSERT INTO users (username, email) VALUES
--   ('tor', 'tor@hjkl.no'),
--   ('extra', 'extra@hjkl.no'),
--   ('ample', 'ample@hjkl.no');