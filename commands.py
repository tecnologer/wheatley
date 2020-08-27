
class Commands:
    # basic commands
    start = "start"
    help = "help"
    cancel = "cancel"

    # configuration
    set_client_id = "setclientid"
    set_secret_id = "setsecretid"

    # users
    add_user = "adduser"
    remove_user = "removeuser"
    get_users = "getusers"
    stream_status = "streamstatus"

    # admins
    add_admin = "addadmin"
    remove_admin = "removeadmin"
    set_master_chat = "setmasterchat"

    def __init__(self):
        super().__init__()

        self.start = Command(self.start, False)
        self.start.set_info('en', "Starts the bot")
        self.start.set_info('es', "Inicia el bot")

        self.help = Command(self.help, False)
        self.help.set_info('en', "Shows the basic info of commands")
        self.help.set_info(
            'es', "Muestra la informacion basica de los comandos")

        self.cancel = Command(self.cancel)
        self.cancel.set_info('en', "Cancels the active command")
        self.cancel.set_info('es', "Cancela el comando activo")

        self.set_client_id = Command(self.set_client_id, False)
        self.set_client_id.set_info('en', "Stores the twitch client id")
        self.set_client_id.set_info('es', "Guarda el client id de twitch")

        self.set_secret_id = Command(self.set_secret_id, False)
        self.set_secret_id.set_info('en', "Stores the twitch client secret")
        self.set_secret_id.set_info('es', "Guarda el client secret de twitch")

        self.add_user = Command(self.add_user)
        self.add_user.set_info(
            'en', "Adds a new user(s) to the list to monitor its stream status. Use users separated by space to add multiple.")
        self.add_user.set_info(
            'es', "Agrega un nuevo usuario a la lista para monitorear si su stream esta activo. Usa usuarios separados por espacio si quieres agregar varios")

        self.remove_user = Command(self.remove_user)
        self.remove_user.set_info(
            'en', "Removes a user(s) from the list to monitor its stream status. Use users separated by space to add multiple.")
        self.remove_user.set_info(
            'es', "Quita al usuario de la lista para no monitorear si su stream esta activo. Usa usuarios separados por espacio si quieres quitar varios")

        self.get_users = Command(self.get_users)
        self.get_users.set_info(
            'en', "Removes a user(s) from the list to monitor its stream status. Use users separated by space to add multiple.")
        self.get_users.set_info(
            'es', "Quita al usuario de la lista para no monitorear si su stream esta activo. Usa usuarios separados por espacio si quieres quitar varios")

        self.stream_status = Command(self.stream_status)
        self.stream_status.set_info(
            'en', "Returns the status for selected user(s). Use separated users by space for multiple users.")
        self.stream_status.set_info(
            'es', "Regresa el estado del stream para el usuario especificado. Usa usuarios separados por espacio si quieres obtener el estatus de varios")

        self.add_admin = Command(self.add_admin, False)
        self.add_admin.set_info(
            'en', "Adds a new telegram's user to whitelist (permissions). Use users separated by space to add multiple.")
        self.add_admin.set_info(
            'es', "Agrega un nuevo usuario de telegram a la lista de administradores. Usa usuarios separados para agregar varios")

        self.remove_admin = Command(self.remove_admin, False)
        self.remove_admin.set_info(
            'en', "Removes a telegram's user from whitelist (permissions). Use users separated by space to add multiple.")
        self.remove_admin.set_info(
            'es', "Quita el usuario de telegram de la lista de administradores. Usa usuarios separados para quitar varios")

        self.set_master_chat = Command(self.set_master_chat, False)
        self.set_master_chat.set_info(
            'en', "Marks this chat as master to recive notifications (debugging purpose).")
        self.set_master_chat.set_info(
            'es', "Marac este canal como maestro para recibir notificaciones con porposito de debugging")


class Command:
    cmd = None
    info = None
    show_in_help = False
    handle = None

    def __init__(self, cmd, in_help=True):
        super().__init__()

        self.info = {}
        self.cmd = cmd
        self.show_in_help = in_help

    def set_info(self, lang, info):
        self.info[lang] = info

    def get_info(self, lang):
        if lang not in self.info:
            lang = 'en'

        return self.info[lang]

    def set_handle(self, handle):
        self.handle = handle
        return self
