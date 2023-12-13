WITH answers_summary AS (
    SELECT a.*, COUNT(*) AS count, AVG(f.rating) AS rating
        FROM "answers" a
        LEFT JOIN "feedback" f
            ON a.id = f.answer
        GROUP BY a.id
), question_summary AS (
    SELECT q.*, COUNT(*) AS count, AVG(a.rating) AS rating
        FROM "questions" q
        LEFT JOIN "answers_summary" a
            ON q.id = a.question
        GROUP BY q.id
), assignment_summary AS (
    SELECT a.*, COUNT(*) AS count, AVG(q.rating) AS rating
        FROM "assignments" a
        LEFT JOIN "question_summary" q
            ON a.id = q.assignment
        GROUP BY a.id
), course_summary AS (
    SELECT c.*, COUNT(*) AS count, AVG(a.rating)
        FROM "courses" c
        LEFT JOIN "assignment_summary" a
            ON c.id = a.course
        GROUP BY c.id
)
SELECT * FROM "course_summary"
