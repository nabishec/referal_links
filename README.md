# Функционал

**POST /auth/register** - авторизация пользователя, возвращает JWT токен

**POST /auth/login** - аутентификация пользователя, возвращает JWT токен

**POST /referrals** - создание реферального кода

**DELETE /referrals** - удаление реферального кода

**GET /referrals/by-email** - получение реферального кода по адресу реферера

**POST /referrals/register** - регистрация по реферальному коду

**GET /referrals/{referrer_id}/referrals** - получение информации о рефералах

**GET /docs** - предоставляет доступ к документации Swagger