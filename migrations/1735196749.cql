CREATE TABLE auth_providers (
    auth_provider_id UUID,
    user_id UUID,
    provider TEXT,
    provider_user_id TEXT,
    family_name TEXT,
    given_name TEXT,
    email TEXT,
    picture TEXT,
    verified_email BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((user_id), provider)
)
