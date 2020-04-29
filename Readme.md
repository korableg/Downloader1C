Консольное приложение для подготовки библиотеки дистрибутива 1С (полезно использовать франчам, свежие обновления всегда под рукой)

Лицензия MIT

Поддерживает следующие флаги:

Обязательные:
-login: Логин на сервер 1С (к сайту releases.1c.ru)
-password: Пароль на сервер 1С (к сайту releases.1c.ru)

Необязательные:
-path: Путь куда складывать скаченные дистрибутивы (по умолчанию "./")
-startdate: Минимальная дата релиза
-nicks: Имена приложений, разделенные запятой (например "platform83, EnterpriseERP20"), подсмотреть можно в адресе, ссылки имею вид например https://releases.1c.ru/project/EnterpriseERP20
-log: Путь к лог файлу, в который сохраняются ошибки, по умолчанию ("./downloader.log")
-h: Справка

Скомпилировано только под архитектуру win64, если нужно компилировать под другие напишите, сделаем.

Console application for preparing the distribution library 1C
Supports the following flags:

Required:
-login: Login to the 1C server (to the site releases.1c.ru)
-password: Password for 1C server (to the site releases.1c.ru)

Optional:
-path: The path to where to download downloaded distributions (default is "./")
-startdate: Minimum release date
-nicks: Application names separated by a comma (for example, "platform83, EnterpriseERP20") can be peeped at the address, the links look like for example https://releases.1c.ru/project/EnterpriseERP20
-log: The path to the log file to which errors are saved, by default ("./downloader.log")
-h: Help

Compiled only for win64 architecture, if you need to compile for others write, we will do it.
