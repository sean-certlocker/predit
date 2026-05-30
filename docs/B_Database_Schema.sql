-- B. Database Schema (PostgreSQL)

CREATE TYPE market_status AS ENUM (
    'DRAFT', 'MODERATION_PENDING', 'APPROVED', 'OPEN', 'LOCKED', 
    'EVENT_ACTIVE', 'RESOLUTION_PENDING', 'RESOLVED', 'VOID', 'MANUAL_REVIEW', 'SETTLED', 'REFUNDED'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'viewer', -- creator, viewer, admin
    trust_score INT DEFAULT 50,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'live',
    start_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    recording_url TEXT,
    health_score INT DEFAULT 100
);

CREATE TABLE markets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES streams(id),
    creator_id UUID REFERENCES users(id),
    title TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    status market_status DEFAULT 'DRAFT',
    
    event_window_seconds INT NOT NULL,
    betting_closes_before_start_seconds INT NOT NULL,
    resolution_method VARCHAR(50) NOT NULL,
    
    max_stake_per_user DECIMAL(18,2),
    max_pool DECIMAL(18,2),
    
    rules JSONB NOT NULL, -- includes outcomes, qualifying objects, void rules
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    opened_at TIMESTAMP WITH TIME ZONE,
    locked_at TIMESTAMP WITH TIME ZONE,
    event_started_at TIMESTAMP WITH TIME ZONE,
    event_ended_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE market_outcomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    market_id UUID REFERENCES markets(id),
    outcome_label VARCHAR(100) NOT NULL, -- e.g. "YES", "NO"
    pool_amount DECIMAL(18,2) DEFAULT 0
);

CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    market_id UUID REFERENCES markets(id),
    outcome_id UUID REFERENCES market_outcomes(id),
    stake DECIMAL(18,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE wallets (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    balance DECIMAL(18,2) DEFAULT 0 CHECK (balance >= 0)
);

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    amount DECIMAL(18,2) NOT NULL, -- positive or negative
    type VARCHAR(50) NOT NULL, -- deposit, stake, win, fee, refund, withdrawal
    reference_id UUID, -- links to position_id, market_id, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- AI & Audit Tables
CREATE TABLE ai_resolution_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    market_id UUID REFERENCES markets(id),
    status VARCHAR(50),
    confidence DECIMAL(5,4),
    raw_output JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE moderation_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type VARCHAR(50), -- market, stream
    target_id UUID,
    action VARCHAR(50), -- approved, rejected, flagged
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
