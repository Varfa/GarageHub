# GarageHub Data Model

## Основной принцип

Структура системы проектируется с учетом реальной работы автосервисов.

Архитектура не считается окончательной и может изменяться по мере эксплуатации системы и получения обратной связи от пользователей.

Все изменения должны сохранять совместимость с существующими данными либо сопровождаться миграциями базы данных.

Client
- number
- name
- phone
- email
- address
- note
- last_visit_at
Car
- client_id
- vin
- plate_number
- brand
- model
- power_kw
- engine
- mileage
- color
Order / Repair
- client_id
- car_id
- accepted_at
- status
- mechanic_id
- work_types
- works
- parts
- photos_before
- photos_after
- total_price
Warehouse
- part_name
- quantity
- unit
- purchase_price
- sale_price
- location
- note
Reports
- repair report for client
- clients report by date / car brand
- mechanic monthly income report
