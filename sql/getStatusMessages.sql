SELECT 
GREATEST(sum(CASE WHEN status < 'sent' THEN 1 END),0),
GREATEST(sum(CASE WHEN status > 'fail' THEN 1 END),0)
FROM message;