# Telegram-бот weath-bot

Созадно с помощью библиотеки echotron: https://www.weatherapi.com

Для запуска через docker-контейнер, необходимо создать .env файл следующего содержания:

```
TG_TOKEN:
WEATHER_TOKEN:
```
Где `TG_TOKEN` — токен бота выданный Telegram, `WEATHER_TOKEN` — токен API, выданынй https://www.weatherapi.com.

Рализованы вызовы для вывода иформации о текущей погоде и предсказания погоды на ближайшие 3 дня, с помощью inline-клавиатуры.
