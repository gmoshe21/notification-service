INSERT INTO message(date_creation,status,id_notification,id_client) VALUES($1,$2,$3,$4) RETURNING id;