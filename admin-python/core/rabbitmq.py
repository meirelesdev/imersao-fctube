from kombu import Connection

def create_rabbitmq_connection()-> Connection:
    return Connection("amqp://admin:admin@rabbitmq:5672//")