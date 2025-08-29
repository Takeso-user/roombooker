-- +migrate Up
INSERT INTO offices (id, name, timezone) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Main Office', 'America/New_York');

INSERT INTO floors (id, office_id, number, label) VALUES
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 1, 'Floor 1'),
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 2, 'Floor 2'),
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 3, 'Floor 3'),
('550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', 4, 'Floor 4'),
('550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440000', 5, 'Floor 5'),
('550e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440000', 6, 'Floor 6'),
('550e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440000', 7, 'Floor 7'),
('550e8400-e29b-41d4-a716-446655440008', '550e8400-e29b-41d4-a716-446655440000', 8, 'Floor 8');

INSERT INTO rooms (id, floor_id, name, capacity, equipment, has_graph_integration, color) VALUES
('550e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440001', 'Room 101', 4, '{"screen": true, "vc": true}', false, '#FF5733'),
('550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440001', 'Room 102', 6, '{"screen": true, "whiteboard": true}', false, '#33FF57'),
('550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440001', 'Room 103', 4, '{"vc": true}', false, '#3357FF'),
('550e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440001', 'Room 104', 8, '{"screen": true, "vc": true, "whiteboard": true}', false, '#FF33A1'),
('550e8400-e29b-41d4-a716-446655440013', '550e8400-e29b-41d4-a716-446655440001', 'Room 105', 6, '{"screen": true}', false, '#A133FF'),
('550e8400-e29b-41d4-a716-446655440014', '550e8400-e29b-41d4-a716-446655440001', 'Room 106', 4, '{"whiteboard": true}', false, '#33FFA1');

-- Add more rooms for other floors similarly, but for brevity, only floor 1

INSERT INTO users (id, email, display_name, role, timezone) VALUES
('550e8400-e29b-41d4-a716-446655440015', 'admin@example.com', 'Admin User', 'admin', 'America/New_York'),
('550e8400-e29b-41d4-a716-446655440016', 'user@example.com', 'Regular User', 'user', 'America/New_York');

INSERT INTO booking_rules (id, office_id, workday_start, workday_end, max_duration, min_lead_time, buffer_before, buffer_after, allow_recurring, timezone) VALUES
('550e8400-e29b-41d4-a716-446655440017', '550e8400-e29b-41d4-a716-446655440000', '09:00:00', '18:00:00', '4 hours', '30 minutes', '15 minutes', '15 minutes', true, 'America/New_York');

-- +migrate Down
DELETE FROM booking_rules WHERE id = '550e8400-e29b-41d4-a716-446655440017';
DELETE FROM users WHERE id IN ('550e8400-e29b-41d4-a716-446655440015', '550e8400-e29b-41d4-a716-446655440016');
DELETE FROM rooms WHERE floor_id = '550e8400-e29b-41d4-a716-446655440001';
DELETE FROM floors WHERE office_id = '550e8400-e29b-41d4-a716-446655440000';
DELETE FROM offices WHERE id = '550e8400-e29b-41d4-a716-446655440000';
