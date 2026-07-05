# GarageHub UI Guide

## Общие правила

- Код интерфейса пишется на английском языке.
- Пользовательский текст может быть на русском, английском или литовском.
- Все переводимые элементы должны иметь `data-i18n`.
- CSS-классы пишутся на английском языке.
- Структура страниц должна быть понятной без комментариев.

---

## Layout

Основной шаблон:

```text
frontend/templates/layout.html
```

Общий layout содержит:

- `header`
- `nav`
- `section.content`

Контент конкретной страницы вставляется через:

```html
{{ template "content" . }}
```

---

## Page Header

Для заголовка страницы используется:

```html
<div class="page-header">
    <h1 data-i18n="title.clients">Клиенты</h1>
    <p data-i18n="message.clients_description">Управление клиентской базой.</p>
</div>
```

---

## Toolbar

Для панели действий страницы:

```html
<div class="clients-toolbar">
    <button class="btn-primary" data-i18n="button.add_client">
        ➕ Добавить клиента
    </button>

    <input
        class="search-input"
        type="text"
        placeholder="Поиск клиента..."
        data-i18n-placeholder="placeholder.search_client"
    >
</div>
```

---

## Buttons

Основная кнопка:

```html
<button class="btn-primary" data-i18n="button.save">Сохранить</button>
```

Обычная кнопка:

```html
<button class="btn-secondary" data-i18n="button.cancel">Отмена</button>
```

Опасное действие:

```html
<button class="btn-danger" data-i18n="button.delete">Удалить</button>
```

---

## Tables

Таблица страницы должна иметь понятный класс:

```html
<table class="clients-table">
    <thead>
        <tr>
            <th data-i18n="table.id">ID</th>
            <th data-i18n="table.name">Имя</th>
            <th data-i18n="table.phone">Телефон</th>
            <th data-i18n="table.cars_count">Автомобилей</th>
            <th data-i18n="table.last_visit">Последний визит</th>
            <th data-i18n="table.actions">Действия</th>
        </tr>
    </thead>

    <tbody>
    </tbody>
</table>
```

---

## Translation attributes

Обычный текст:

```html
data-i18n="nav.clients"
```

Placeholder:

```html
data-i18n-placeholder="placeholder.search_client"
```

Title:

```html
data-i18n-title="title.edit_client"
```

ARIA label:

```html
data-i18n-aria-label="aria.edit_client"
```

---

## Translation key groups

- `nav.*` — меню
- `button.*` — кнопки
- `title.*` — заголовки страниц
- `label.*` — подписи полей
- `table.*` — заголовки таблиц
- `message.*` — сообщения пользователю
- `placeholder.*` — placeholder полей
- `error.*` — ошибки
- `aria.*` — accessibility labels

---

## Naming

Хорошо:

```html
class="clients-toolbar"
class="clients-table"
class="search-input"
```

Плохо:

```html
class="block1"
class="table"
class="knopka"
```

---

## Principle

Если элемент относится к конкретной странице, класс содержит имя страницы:

```text
clients-toolbar
cars-table
orders-filter
```

Если элемент общий для всего проекта, класс общий:

```text
btn-primary
btn-secondary
search-input
page-header
content
```
