# GarageHub Roadmap

Версия: 1.0

---

# Статус проекта

🟢 Проектирование — завершено

✅ Техническое задание

✅ Архитектура

✅ Структура проекта

✅ Документация

✅ Среда разработки

---

# Этап 1 — Backend MVP

## Sprint 1 — Infrastructure

Статус: 🟡 В процессе

Задачи:

- [x] Настроить конфигурацию проекта
- [x] Подключить PostgreSQL
- [x] Реализовать Connection Pool
- [x] Реализовать Health endpoint
- [ ] Реализовать единый Response package
- [ ] Настроить Logger
- [x] Первый успешный запуск Backend

Результат:

Рабочий backend без бизнес-логики.

---

## Sprint 1.5 — Frontend Foundation

Статус: 🟡 В процессе

Задачи:

- [x] Создать базовую структуру HTML
- [x] Подключить Go Templates
- [x] Настроить раздачу статических файлов
- [x] Подключить CSS
- [x] Реализовать Header
- [x] Реализовать Sidebar
- [x] Создать первый Dashboard
- [ ] Реализовать адаптивную верстку
- [ ] Подключить Font Awesome
- [ ] Реализовать переключение языков

Результат:

Рабочий интерфейс CRM с базовой навигацией.

---

## Sprint 2 — Authentication

Статус: ⏳

Задачи:

- [ ] Users
- [ ] Login
- [ ] Password Hashing
- [ ] JWT
- [ ] Middleware Authorization

---

## Sprint 3 — Clients

Статус: ⏳

Задачи:

- [ ] CRUD Clients
- [ ] Draft Clients
- [ ] Validation
- [ ] Search

---

## Sprint 4 — Cars

Статус: ⏳

Задачи:

- [ ] CRUD Cars
- [ ] Car Details
- [ ] VIN Validation

---

## Sprint 5 — Repairs

Статус: ⏳

Задачи:

- [ ] CRUD Repairs
- [ ] Repair Status
- [ ] Repair History

---

## Sprint 6 — Reports

Статус: ⏳

Задачи:

- [ ] Reports
- [ ] Statistics
- [ ] Filters

---

## Sprint 7 — Notifications

Статус: ⏳

Задачи:

- [ ] Telegram
- [ ] WhatsApp
- [ ] SMS

---

## Sprint 8 — Admin Panel API

Статус: ⏳

---

## Sprint 9 — Testing

Статус: ⏳

---

## Sprint 10 — Production

Статус: ⏳

Задачи:

- [ ] Docker
- [ ] Docker Compose
- [ ] CI/CD
- [ ] Deploy
- [ ] Nginx
- [ ] HTTPS

---

# Этап 2 — Frontend

- [x] HTML5
- [x] CSS3
- [ ] JavaScript
- [x] Go Templates
- [ ] Responsive Design
- [ ] Font Awesome
- [ ] Dark Theme

---

# Этап 3 — Release

GarageHub v1.0

---

# Архитектурные принципы

- Документировать важные архитектурные решения.
- Код полностью понятен разработчику без копирования.
- Backend и Frontend развиваются параллельно.
- Все пользовательские тексты поддерживают мультиязычность.
- Для переводов используется `data-i18n`.
- GarageHub разрабатывается как коммерческий продукт, пригодный для использования на реальном СТО.

---

# Главная цель

GarageHub должен стать полноценной системой управления автосервисом, которую можно использовать в реальной работе и не стыдно показать работодателю или заказчику.
