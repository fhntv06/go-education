# go-education

Структура проекта:

<pre><code>project/
│
├── backend/
│   └── main.go
│
└── frontend/
    ├── index.html
    └── assets/
        ├── app.js
        └── app.css
</code></pre>


## Backend
<pre><code>backend/
│
├── main.go           # Основной файл сервера
├── handlers/         # Папка с обработчиками HTTP-запросов
│   └── ...           # Файлы обработчиков
├── models/           # Модели данных
│   └── ...           # Файлы моделей
└── db/               # Конфигурация базы данных
└── config.json   # Настройки подключения к базе данных
</code></pre>

## Frontend
<pre><code>frontend/
│
├── index.html       # Главная страница
├── assets/          # Статические ресурсы
│   ├── styles/      # Папка с таблицами стилей
│   │   └── app.css  # Основные стили
│   └── scripts/     # Скрипты JavaScript
│       └── app.js   # Главный скрипт
└── components/      # Компоненты React/Vue/Angular (если используются)
└── ...          # Файлы компонентов
</code></pre>


