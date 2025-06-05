CREATE OR REPLACE PROCEDURE unblock_user(
    p_user_id INT DEFAULT NULL,
    p_name VARCHAR DEFAULT NULL,
    p_phone VARCHAR[] DEFAULT NULL,
    p_email VARCHAR DEFAULT NULL,
    p_status INT DEFAULT NULL,  
    p_block_reason_code INT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_user_id IS NOT NULL THEN
        UPDATE users
        SET block_reason_code = 0,
            block_reason = NULL,
            status = 1,
            updated_at = NOW()
        WHERE id = p_user_id;
    ELSE
        UPDATE users
        SET block_reason_code = 0,
            block_reason = NULL,
            status = 1,
            updated_at = NOW()
        WHERE 
            (p_name IS NULL OR name ILIKE '%' || p_name || '%') AND
            (p_phone IS NULL OR phone = ANY (p_phone)) AND
            (p_email IS NULL OR email = p_email) AND
            (p_status IS NULL OR status = p_status) AND
            (p_block_reason_code IS NULL OR block_reason_code = p_block_reason_code);
    END IF;
END;
$$;
