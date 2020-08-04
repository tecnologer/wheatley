
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
    get_users = "getusers",

    # admins
    add_admin = "addadmin"
    remove_admin = "removeadmin"
    set_master_chat = "setmasterchat"

    def __init__(self):
        super().__init__()
