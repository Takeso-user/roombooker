-- +migrate Down
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS holidays;
DROP TABLE IF EXISTS booking_rules;
DROP TABLE IF EXISTS booking_participants;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS floors;
DROP TABLE IF EXISTS offices;
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS users;
