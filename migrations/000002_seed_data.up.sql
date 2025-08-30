-- +migrate Up

INSERT INTO offices (id, name, timezone) VALUES
('office-1', 'Main Office', 'America/New_York');

INSERT INTO floors (id, office_id, number, label) VALUES
('floor-1-1', 'office-1', 1, 'Floor 1'),
('floor-1-2', 'office-1', 2, 'Floor 2'),
('floor-1-3', 'office-1', 3, 'Floor 3'),
('floor-1-4', 'office-1', 4, 'Floor 4'),
('floor-1-5', 'office-1', 5, 'Floor 5'),
('floor-1-6', 'office-1', 6, 'Floor 6'),
('floor-1-7', 'office-1', 7, 'Floor 7'),
('floor-1-8', 'office-1', 8, 'Floor 8');

INSERT INTO rooms (id, floor_id, name, capacity, equipment, has_graph_integration, color) VALUES
('room-101', 'floor-1-1', 'Room 101', 4, '{"screen": true, "vc": true}', 0, '#FF5733'),
('room-102', 'floor-1-1', 'Room 102', 6, '{"screen": true, "whiteboard": true}', 0, '#33FF57'),
('room-103', 'floor-1-1', 'Room 103', 4, '{"vc": true}', 0, '#3357FF'),
('room-104', 'floor-1-1', 'Room 104', 8, '{"screen": true, "vc": true, "whiteboard": true}', 0, '#FF33A1'),
('room-105', 'floor-1-1', 'Room 105', 6, '{"screen": true}', 0, '#A133FF'),
('room-106', 'floor-1-1', 'Room 106', 4, '{"whiteboard": true}', 0, '#33FFA1');

INSERT INTO users (id, email, display_name, role, timezone) VALUES
('user-admin', 'admin@example.com', 'Admin User', 'admin', 'America/New_York'),
('user-regular', 'user@example.com', 'Regular User', 'user', 'America/New_York');

INSERT INTO booking_rules (id, office_id, workday_start, workday_end, max_duration, min_lead_time, buffer_before, buffer_after, allow_recurring, timezone) VALUES
('rule-1', 'office-1', '09:00:00', '18:00:00', '4 hours', '30 minutes', '15 minutes', '15 minutes', 1, 'America/New_York');

-- +migrate Down
DELETE FROM booking_rules WHERE id = 'rule-1';
DELETE FROM users WHERE id IN ('user-admin', 'user-regular');
DELETE FROM rooms WHERE floor_id LIKE 'floor-1-%';
DELETE FROM floors WHERE office_id = 'office-1';
DELETE FROM offices WHERE id = 'office-1';
