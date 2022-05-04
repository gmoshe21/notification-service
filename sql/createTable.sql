CREATE TABLE IF NOT EXISTS notification (
    id    uuid NOT NULL,
    start_data timestamp with time zone NOT NULL,
    text_message varchar(255) NOT NULL,
    filter_clients  varchar(255) NOT NULL,
    finish_data timestamp with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS client (
    id    uuid NOT NULL,
    phone_number    bigserial NOT NULL,
    mobile_operator_code    integer NOT NULL,
    teg    varchar(10) NOT NULL,
    time_zone    time with time zone
);

CREATE TABLE IF NOT EXISTS message (
    id    bigserial PRIMARY KEY,
    date_creation timestamp NOT NULL,
    status varchar(255) NOT NULL,
    id_notification uuid NOT NULL,
    id_client uuid NOT NULL
);
