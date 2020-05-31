<h3>Консольное приложение для подготовки библиотеки дистрибутива 1С (полезно использовать франчам, свежие обновления всегда под рукой)</h3>

Поддерживает следующие флаги:

<h4>Обязательные:</h4>
<b>-login</b>: Логин на сервер 1С (к сайту releases.1c.ru)<br>
<b>-password</b>: Пароль на сервер 1С (к сайту releases.1c.ru)

<h4>Необязательные:</h4>
<b>-path</b>: Путь куда складывать скаченные дистрибутивы (по умолчанию "./")<br>
<b>-startdate</b>: Минимальная дата релиза <br>
<b>-nicks</b>: Имена приложений, разделенные запятой (например "platform83, EnterpriseERP20"), подсмотреть можно в адресе, ссылки имею вид например https://releases.1c.ru/project/EnterpriseERP20 <br>
<b>-log</b>: Путь к лог файлу, в который сохраняются ошибки, по умолчанию ("./downloader.log")<br>
<b>-h</b>: Справка<br>

<h4>Команды сервиса:</h4>
<b>install</b>: Установить сервис<br>
<b>remove</b>: Удалить сервис<br>
<b>start</b>: Запустить сервис<br>
<b>stop</b>: Остановить сервис<br>
<b>pause</b>: Поставить сервис на паузу (Активный процесс скачивания будет работать пока не завершится)<br>
<b>continue</b>: Продолжить работу (после паузы)<br>

<h4>Дополнительно для сервиса:</h4>
<b>-instance</b>: Название сервиса (По умолчанию Downloader1C) (на случай если требуется развернуть несколько)
<br><br>
Скомпилировано только под архитектуру win64, если нужно компилировать под другие напишите, сделаем.

<h3>Console application for preparing the distribution library 1C</h3>
Supports the following flags:

<h4>Required:</h4>
<b>-login</b>: Login to the 1C server (to the site releases.1c.ru)<br>
<b>-password</b>: Password for 1C server (to the site releases.1c.ru)<br>

<h4>Optional:</h4>
<b>-path</b>: The path to where to download downloaded distributions (default is "./")<br>
<b>-startdate</b>: Minimum release date<br>
<b>-nicks</b>: Application names separated by a comma (for example, "platform83, EnterpriseERP20") can be peeped at the address, the links look like for example https://releases.1c.ru/project/EnterpriseERP20 <br>
<b>-log</b>: The path to the log file to which errors are saved, by default ("./downloader.log")<br>
<b>-h</b>: Help<br>
<br><br>
Compiled only for win64 architecture, if you need to compile for others write, we will do it.

<h4>For service:</h4>
<b>install</b>: Install service<br>
<b>remove</b>: Remove service<br>
<b>start</b>: Start service<br>
<b>stop</b>: Stop service<br>
<b>pause</b>: Pause service<br>
<b>continue</b>: Continue service<br>

<h4>Additional argunments:</h4>
<b>-instance</b>: Service name (default Downloader1C) (for more instances)
