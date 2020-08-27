
class Language:
    def __init__(self, update):
        super().__init__()
        self.code = update.effective_user.language_code

        if self.code is None or self.code not in self.__languages:
            self.code = 'en'

    code = 'en'
    __languages = {
        "es": {
            "bot_description": "Bot para recibir notificaciones cuando un streamer de twitch entre en linea",
            "help_header": "Comandos disponibles:",
            "client_id_request": "Favor de escribir el Client ID de twitch:",
            "client_id_success": "Client Id configurado",
            "client_secret_request": "Favor de escribir el Client secret de twitch:",
            "client_secret_success": "Client secret configurado",
            "add_user_success": "Ahora te notificare en este chat cuando {0} inicie su stream",
            "add_user_limit_reached": "El limite de usuarios ha sido alcanzado, favor de contactar alguno de nuestros administradores: {0}",
            "add_user_request_name": "Escribe el nombre de usuario de twitch:",
            "remove_use_success": "Las notificaciones de {0} han sido apagadas",
            "remove_user_not_found": "El usuario(s) {0} no  esta configurado(s)",
            "cancel_success": "Listo!, Todo esta cancelado.",
            "add_admin_success": "{0} es admin 😉",
            "add_admin_no_specified": "El nombre de usuario para admin es requerido. Usa: /{0} <@telegram_username>",
            "remove_admin_no_specified": "El nombre de usuario para admin es requerido. use: /{0} <@telegram_username>",
            "remove_admin_success": "{0} ya no es admin 😢",
            "chat_master_no_change": "Nada que hacer",
            "chat_master_success": "Ahora te notificare aqui cualquier cosa rara que pase",
            "chat_master_error": "Ocurrio una excepcion",
            "get_users_no_users": "No hay usuarios de twitch configurados en este chat",
            "get_users_header": "*Los siguientes usuarios estan configurados en este chat:*\n\n",
            "get_users_header_invalid_users": "{0}\n*Los siguientes nombres no son usuarios validos en twitch:*\n\n",
            "stream_status_error": "El nombre de usuario es requerido. /{0} <username>[ <username>]",
            "is_allowed_fail": "Permiso requerido",
            "config_assist_client_id": "Antes de continuar es necesario configurar el client id de twitch",
            "config_assist_client_secret": "Tambien es necesario configurar el client secret de twitch.",
            "config_assist_greeting": "Hola, te voy a asistir para configurar todo lo necesario",
            "config_assist_greeting": "Felicidades!, todo esta configurado."
        },
        "en": {
            "bot_description": "Bot to receive notifications when a twitch's user starts his stream",
            "help_header": "Available commands:",
            "client_id_request": "Please type the twitch Client ID:",
            "client_id_success": "Client Id configured",
            "client_secret_request": "Please type the twitch Client secret:",
            "client_secret_success": "Client secret configured",
            "add_user_success": "Now I'll notify you in this chat when {0} is streaming",
            "add_user_limit_reached": "The limit of users has been reached, please contact any of this administrators: {0}",
            "add_user_request_name": "Type the username of User's twitch:",
            "remove_user_success": "The notifications for {0} are turned off",
            "remove_user_not_found": "The user {0} is not configured",
            "remove_user_request_name": "Type the username of User's twitch:",
            "cancel_success": "Done!, everything is canceled.",
            "add_admin_success": "{0} is admin 😉",
            "add_admin_no_specified": "The name of admin(s) is required. use: /{0} <telegram_username>",
            "remove_admin_no_specified": "The name of admin(s) is required. use: /{0} <@telegram_username>",
            "remove_admin_success": "{0} are no longer admin 😢",
            "chat_master_no_change": "Noting to chance",
            "chat_master_success": "Now I'll notify you here if any weird happens",
            "chat_master_error": "an exception occurred",
            "get_users_no_users": "There are not users twitch configured in this chat",
            "get_users_header": "*The following users twitch are registered in this chat:*\n\n",
            "get_users_header_invalid_users": "{0}\n*The following users twitch are invalid:*\n\n",
            "stream_status_error": "The user's name is required. /{0} <username>[ <username>]",
            "is_allowed_fail": "Permissions required",
            "config_assist_client_id": "Before continue, We need set the twitch client id.",
            "config_assist_client_secret": "Also, We need set the twitch client secret.",
            "config_assist_greeting": "Hello, I'll assist you to configure everything",
            "config_assist_greeting": "Congratulations!, everything is configured."
        }
    }

    def get(self, key):
        if key not in self.__languages[self.code]:
            return ""

        return self.__languages[self.code][key]
