# Weather Subscription Service

### ✅ Сервіс розгорнуто на VPS
[Swagger](http://vps71203.hyperhost.name:8080/swagger/index.html)

### 🧱 Архітектура

- `internal/handlers` — HTTP хендлери
- `internal/service` — бізнес-логіка
- `internal/storage` — інтерфейс до бази даних (PostgreSQL)
- `internal/clients/weatherapi` — клієнт до [weatherapi.com](https://www.weatherapi.com/)

---

## 📦 Технології

- sqlc для генерації SQL-коду
- HTTP клієнт `resty` для запитів до weatherapi
- Docker + docker-compose

---

## 🚀 Потенційні покращення

### 🔐 Генерація токенів

Варто перейти на тимчасові токени з обмеженим часом дії (наприклад, через JWT або `exp` в базі даних).

### 📩 Надсилання повідомлень

Потрібно реалізувати окремий **email-клієнт**, який буде надсилати прогноз на email користувачам згідно з частотою підписки.

### ⚡️ Кешування погоди

Щоб не дублювати запити до weatherapi.com, можна реалізувати **кешування погоди по містах** на деякий час (наприклад, 10 хв) через Redis:

- ключ: `weather:<city>`
- значення: JSON прогнозу
- TTL: 10–15 хвилин

Це суттєво зменшить кількість запитів до зовнішнього API та пришвидшить `/weather`.

### 👻 Додати .env в .gitignore

Звісно не світити .env в реальному проєкті 

---

## 🧪 Запуск
make up

