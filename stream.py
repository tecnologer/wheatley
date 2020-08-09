import games


def get_from_response(twitch_data):
    stream = Stream()
    stream.user_name = twitch_data['user_name'].lower()
    stream.title = twitch_data['title']
    stream.viewer_count = twitch_data['viewer_count']
    stream.game_id = twitch_data['game_id']
    stream.game_name = games.get_game(stream.game_id)
    stream.started_at = twitch_data['started_at']

    return stream


class Stream:
    user_name = None
    title = None
    viewer_count = 0
    game_id = None
    game_name = None
    started_at = None

    def __init__(self):
        super().__init__()
