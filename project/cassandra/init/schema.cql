CREATE KEYSPACE IF NOT EXISTS chargebacks
    WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1
        };

USE chargebacks;

CREATE TABLE IF NOT EXISTS chargebacks_by_user_transaction (
    user_id TEXT,
    transaction_id TEXT,
    status TEXT,
    reason TEXT,
    file_id UUID,
    sent_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((user_id, transaction_id))
);