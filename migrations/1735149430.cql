INSERT INTO roles (role_id, role_name, created_at, updated_at) 
VALUES (uuid(), 'ADMIN', toTimestamp(now()), toTimestamp(now()));

INSERT INTO roles (role_id, role_name, created_at, updated_at) 
VALUES (uuid(), 'USER', toTimestamp(now()), toTimestamp(now()));

INSERT INTO roles (role_id, role_name, created_at, updated_at) 
VALUES (uuid(), 'PARTNER', toTimestamp(now()), toTimestamp(now()));