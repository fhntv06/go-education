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

## Корректное написание кода
- https://go.dev/talks/2014/names.slide#1

<hr>

## Что сделано по Go!
<table>
    <tr>
        <th>Дата</th>
        <th>Изучил / Сделано </th>
    </tr>
    <tr>
        <td>
            <ul>
                <li>13.01.2025</li>
            </ul>
        </td>
        <td>
            <ul>
                <li>Цикл for</li>
                <li><a href="https://go.dev/talks/2014/names.slide#1">Стилистика кода</a></li>
                <li><a href="https://go.dev/doc/modules/managing-source">Организация модулей</a></li>
                <li><a href="https://go.dev/doc/modules/layout">Организация модуля Go</a></li>
                <li>Ресерчил информацию по пет-проекту</li>
                <li>Подключение к PostgreSQL</li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>
            <ul>
                <li>14.01.2025</li>
            </ul>
        </td>
        <td>
            <ul>
                <li>На работе чуть-чуть туториала</li>
                <li>Сделал структуру проекта AI-generator-image</li>
                <li>Нашел API для генерации изображений (<a href="docs.getimg.ai">Docs getimg</a>)</li>
                <li>Создал новые пакеты для работы в AI-generator-image</li>
            </ul>
        </td>
    </tr>
</table>

<hr>