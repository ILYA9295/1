# 🚀 Anti-Hate Telegram Bot MVP

## 📌 Возможности
- Автоматически анализирует сообщения в чате.
- Удаляет токсичные, матерные или оскорбительные сообщения.
- Сохраняет все сообщения в PostgreSQL (с флагом `deleted`).
- Поднимает мини-админку (`/messages`) для просмотра истории.

---

## ⚙️ Запуск

### 1. Склонировать проект
```sh
git clone <repo>
cd project

#2. Настроить .env
TELEGRAM_BOT_TOKEN=ваш_токен_бота
POSTGRES_USER=bot_user
POSTGRES_PASSWORD=bot_pass
POSTGRES_DB=bot_db
POSTGRES_HOST=db
POSTGRES_PORT=5432
APP_PORT=8080

#3. Сгенерировать gRPC код
Python:
cd nlp_service
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. nlp.proto

Go (PowerShell):
protoc --go_out=. --go-grpc_out=. nlp_service/nlp.proto

#4. Запустить Docker
docker-compose up --build


app → Go бот

db → PostgreSQL

nlp_service → Python gRPC

🤖 Подключение бота
Создайте бота через @BotFather
Вставьте токен в .env.
Добавьте бота в свой чат.
Назначьте его администратором с правом удалять сообщения.

📊 Админка
Здоровье: http://localhost:8080/health
Сообщения: http://localhost:8080/messages

✅ MVP готов
Чистые сообщения → сохраняются в БД + подтверждаются ботом.
Токсичные → удаляются из чата + помечаются в БД (deleted=true).
История доступна через /messages.

3. Генерация gRPC файлов
📍 Установка protoc (Windows)

Скачай Protocol Buffers (protoc)

→ выбери protoc-<версия>-win64.zip.

Распакуй, добавь путь к bin/protoc.exe в PATH.
Проверка:

protoc --version

📍 Python

Установи пакеты:

pip install grpcio grpcio-tools


Сгенерируй:

cd nlp_service
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. nlp.proto


→ появятся файлы nlp_pb2.py и nlp_pb2_grpc.py.

📍 Go

Установи плагины:

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


Добавь $GOPATH/bin в PATH. Проверка:

protoc-gen-go --version
protoc-gen-go-grpc --version


Сгенерируй:

protoc --go_out=. --go-grpc_out=. nlp_service/nlp.proto


→ появятся nlp_service/nlp.pb.go.