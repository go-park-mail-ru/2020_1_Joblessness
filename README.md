![Docker build](https://github.com/go-park-mail-ru/2020_1_Joblessness/workflows/Docker%20build/badge.svg)
# Joblessness
### Репозиторий бэкэнда команды Joblessness
### https://hahao.ru/

### API

| Method name | Path | Description |
| ----------- | ---- | ----------- |
| POST | /users/login | Войти |
| POST | /users/check | Проверить валидность сессии |
| POST | /users/logout | Выйти |
| POST | /users | Создать нового пользователя |
| POST | /organizations | Создать новую организацию |
| PUT | /summaries/{summary_id:[0-9]+}/response | Ответить на резюме пользователя |
| ---- | /chat | Войти в чат |
| GET | /chat/conversation/{user_id:[0-9]+} | История сообщений |
| GET | /chat/conversation | Активные диалоги |
| GET | /recommendation | Рекомендованные вакансии |
| GET | /search | Поиск |
| POST | /summaries | Создать резюме |
| GET | /summaries | Получить список резюме |
| GET | /summaries/{summary_id:[0-9]+} | Получить резюме |
| DELETE | /summaries/{summary_id:[0-9]+} | Получить удалить резюме |
| GET | /summaries/{summary_id:[0-9]+}/print | Распечатать резюме в pdf |
| GET | /users/{user_id:[0-9]+}/summaries | Получить список резюме пользователя |
| POST | /vacancies/{vacancy_id:[0-9]+}/response | Откликнуться на вакансию |
| GET | /organizations/{user_id:[0-9]+}/response | Получить список откликов на вакансии организации |
| GET | /users/{user_id:[0-9]+}/response | Получить список откликов на резюме пользователя |
| POST | /summaries/{summary_id:[0-9]+}/mail | Отправить резюме пользователя по почте |
| POST | /vacancies | Создать вакансию |
| GET | /vacancies | Получить список вакансий |
| GET | /vacancies/{vacancy_id:[0-9]+} | Получить вакансию |
| DELETE | /vacancies/{vacancy_id:[0-9]+} | Удалить вакансию |
| GET | /organizations/{organization_id:[0-9]+}/vacancies | Получить список вакансий организации |
| GET | /users/{user_id:[0-9]+} | Получить информацию о пользователе |
| PUT | /users/{user_id:[0-9]+} | Изменить информацию о пользователе |
| GET | /organizations/{user_id:[0-9]+} | Получить информацию об организации |
| PUT | /organizations/{user_id:[0-9]+} | Изменить информацию об организации |
| GET | /organizations | Получить список организаций |
| POST | /users/{user_id:[0-9]+}/avatar | Изменить аватарку пользователя |
| POST | /users/{user_id:[0-9]+}/like | Добавить/удалить пользователя в избранное |
| GET | /users/{user_id:[0-9]+}/like | Проверить нахождение пользователя в избранном |
| GET | /users/{user_id:[0-9]+}/favorite | Получить список избранного пользователя |

API в Swagger: https://app.swaggerhub.com/apis/maratishimbaev/hh.ru/1.0.0

### Авторы

* **Ишимбаев Марат** - [@maratishimbaev](https://github.com/GavrilyukAG)
* **Куклин Сергей** - [@huvalk](https://github.com/RusPatrick)
* **Балицкий Михаил** - [@mikstime](https://github.com/makdenis)



### Менторы

Backend: Артем Бакулев

Frontend: Нозим Юнусов

### Репозиторий Frontend: [github.com](https://github.com/frontend-park-mail-ru/2020_1_Joblessness)