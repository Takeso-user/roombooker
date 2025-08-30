CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    display_name TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    auth_provider TEXT,
    password_hash TEXT,
    mfa_enabled INTEGER DEFAULT 0,
    timezone TEXT DEFAULT 'UTC',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE offices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    timezone TEXT DEFAULT 'UTC'
);

CREATE TABLE floors (
    id TEXT PRIMARY KEY,
    office_id TEXT REFERENCES offices(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    label TEXT
);

CREATE TABLE rooms (
    id TEXT PRIMARY KEY,
    floor_id TEXT REFERENCES floors(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    capacity INTEGER NOT NULL,
    equipment TEXT,
    has_graph_integration INTEGER DEFAULT 0,
    graph_resource_id TEXT,
    color TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
    id TEXT PRIMARY KEY,
    room_id TEXT REFERENCES rooms(id) ON DELETE CASCADE,
    created_by TEXT REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    starts_at_utc DATETIME NOT NULL,
    ends_at_utc DATETIME NOT NULL,
    rrule TEXT,
    status TEXT DEFAULT 'active',
    external_event_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE booking_participants (
    id TEXT PRIMARY KEY,
    booking_id TEXT REFERENCES bookings(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    is_required INTEGER DEFAULT 1,
    response_status TEXT
);

CREATE TABLE booking_rules (
    id TEXT PRIMARY KEY,
    office_id TEXT REFERENCES offices(id) ON DELETE CASCADE,
    workday_start TEXT NOT NULL,
    workday_end TEXT NOT NULL,
    max_duration TEXT NOT NULL,
    min_lead_time TEXT NOT NULL,
    buffer_before TEXT DEFAULT '0 minutes',
    buffer_after TEXT DEFAULT '0 minutes',
    allow_recurring INTEGER DEFAULT 0,
    timezone TEXT DEFAULT 'UTC',
    effective_from DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE holidays (
    id TEXT PRIMARY KEY,
    office_id TEXT REFERENCES offices(id) ON DELETE CASCADE,
    date TEXT NOT NULL,
    description TEXT
);

CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    actor_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    payload_json TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_room_starts ON bookings(room_id, starts_at_utc);
CREATE INDEX idx_rooms_floor ON rooms(floor_id);
CREATE INDEX idx_audit_actor ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
