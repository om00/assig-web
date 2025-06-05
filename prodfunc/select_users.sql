CREATE OR REPLACE FUNCTION select_users(
    name_filter VARCHAR DEFAULT NULL,
    phone_filter VARCHAR[] DEFAULT NULL,
    email_filter VARCHAR DEFAULT NULL,
    block_reason_code_filter INT DEFAULT NULL,
    status_filter INT DEFAULT NULL
)
RETURNS TABLE(
    id INT,
    name VARCHAR,
    age INT,
    phone VARCHAR,
    email VARCHAR,
    status INT,
    block_reason TEXT,
    block_reason_code INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id, 
        u.name, 
        u.age, 
        u.phone, 
        u.email, 
        u.status, 
        u.block_reason,        
        u.block_reason_code,   
        u.created_at, 
        u.updated_at
    FROM users u
    WHERE 
        (name_filter IS NULL OR u.name ILIKE '%' || name_filter || '%') AND
        (phone_filter IS NULL OR u.phone = ANY (phone_filter)) AND
        (email_filter IS NULL OR u.email = email_filter) AND
        (block_reason_code_filter IS NULL OR u.block_reason_code = block_reason_code_filter) AND
        (status_filter IS NULL OR u.status = status_filter);
END;
$$ LANGUAGE plpgsql;