# testAvito
Запуск:
1. Склонировать репозиторий
2. Запустить контейнеры сервера и базы (sudo docker compose up -d --build)
3. Запустить скрипт для создания таблиц в базе (bash createTabkes.sh)
4. Swagger доступен по url http://127.0.0.1:8080/api/docs/index.html#/
5. Дернуть в Swagger ручку filldb для заполнения базы тестовыми данными

Принятые решения:
1. Сделал отдельные таблицы tags и features (т.к. в реальном сервисе скорее всего будут необходимые для этих сущностей поля), но оставил только числовое поле
id, т.к. остальное в рамках тестового задания не нужно.
2. Сделал отдельную таблицу для сопоставления тегов и баннеров, т.к. если хранить теги списком в одном поле в таблице баннеров, то операции поиска 
и удаления баннеров по тегу будут сильно замедляться с ростом числа записей в таблицах. При этом в таблице также хранится поле feature_id для быстрого поиска
баннера по тегу и фиче.
3. Добавил swagger для тестирования.
4. Для удобства тестирования добавил функцию filldb - заполняет бд тестовыми данными. Для пользователя и админа генерируются токены, для их просмотра можно 
воспользоваться функцией GetUsers.
5. В запросе get /banner сервер возвращает массив, а не json-объект, т.к. так указано в задании.
6. Оставил конфигурационные данные для подключение к базе открытыми, не стал выносить в github секреты, так как не успел это сделать.
