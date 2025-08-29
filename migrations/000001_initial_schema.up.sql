-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    auth_provider VARCHAR(50),
    password_hash VARCHAR(255),
    mfa_enabled BOOLEAN DEFAULT FALSE,
    timezone VARCHAR(100) DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    raw_profile_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE offices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    timezone VARCHAR(100) DEFAULT 'UTC'
);

CREATE TABLE floors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    office_id UUID REFERENCES offices(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    label VARCHAR(255)
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    floor_id UUID REFERENCES floors(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    capacity INTEGER NOT NULL,
    equipment JSONB,
    has_graph_integration BOOLEAN DEFAULT FALSE,
    graph_resource_id VARCHAR(255),
    color VARCHAR(7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    starts_at_utc TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at_utc TIMESTAMP WITH TIME ZONE NOT NULL,
    rrule TEXT,
    status VARCHAR(50) DEFAULT 'active',
    external_event_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    EXCLUDE (room_id WITH =, tstzrange(starts_at_utc, ends_at_utc) WITH &&) WHERE (status = 'active')
);

CREATE TABLE booking_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    is_required BOOLEAN DEFAULT TRUE,
    response_status VARCHAR(50)
);

CREATE TABLE booking_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    office_id UUID REFERENCES offices(id) ON DELETE CASCADE,
    workday_start TIME NOT NULL,
    workday_end TIME NOT NULL,
    max_duration INTERVAL NOT NULL,
    min_lead_time INTERVAL NOT NULL,
    buffer_before INTERVAL DEFAULT '0 minutes',
    buffer_after INTERVAL DEFAULT '0 minutes',
    allow_recurring BOOLEAN DEFAULT FALSE,
    timezone VARCHAR(100) DEFAULT 'UTC',
    effective_from TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    office_id UUID REFERENCES offices(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    payload_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_bookings_room_starts ON bookings(room_id, starts_at_utc);
CREATE INDEX idx_rooms_floor ON rooms(floor_id);
CREATE INDEX idx_audit_actor ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);

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
