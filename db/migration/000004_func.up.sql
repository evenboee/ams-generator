CREATE OR REPLACE FUNCTION get_assignment_rating(assignment_id_arg int)
RETURNS float AS $$
DECLARE
    _rating float;
BEGIN
    -- Check if entry exists and the cache hasn't expired
    IF NOT EXISTS(
        SELECT 1
        FROM assignment_rating_cache
        WHERE assignment_id = assignment_id_arg AND expires_at > NOW()
    ) THEN
        -- Calculate the average rating and update the cache
        SELECT rating INTO _rating
        FROM assignment_summary
        WHERE id = assignment_id_arg;

        -- Update or insert the cache
        INSERT INTO assignment_rating_cache (assignment_id, avg_rating, expires_at)
        VALUES (assignment_id_arg, _rating, NOW() + INTERVAL '10 minutes')
        ON CONFLICT (assignment_id)
        DO UPDATE SET avg_rating = _rating, expires_at = NOW() + INTERVAL '10 minutes';
    ELSE
        -- Get the rating from the cache
        SELECT avg_rating INTO _rating
        FROM assignment_rating_cache
        WHERE assignment_id = assignment_id_arg;
    END IF;

    RETURN _rating;
END;
$$ LANGUAGE plpgsql;

CREATE VIEW assignment_with_rating AS
    SELECT *, get_assignment_rating(id) AS rating
        FROM assignments;
