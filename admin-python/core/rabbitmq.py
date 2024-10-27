from kombu import Connection
from django.conf import settings

def create_rabbitmq_connection()-> Connection:
    return Connection(settings.RABBITMQ_URL)