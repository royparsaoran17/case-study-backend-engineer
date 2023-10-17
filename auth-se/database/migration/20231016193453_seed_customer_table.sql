-- +goose Up
INSERT INTO customers (id, name, phone, password, role_id, created_at, updated_at, deleted_at)
VALUES ('a05b5ac4-53c4-11ee-8c99-0242ac120002', 'Roy Parsaoran', '+6281809134100', 'b1a22367af5cb7709d7f7a9211c777c176fdd807', '84756252-dadd-46c4-836f-608a1e8d33ce', '2023-10-09 19:18:05.000000',
        '2023-10-09 19:18:05.000000', null); -- password : JuloTest123

-- +goose Down
DELETE
FROM customers
WHERE id = 'a05b5ac4-53c4-11ee-8c99-0242ac120002';

