FROM python:3.12.6-alpine

WORKDIR /usr/src/app

# Instalar pipenv
RUN pip install --upgrade pip && pip install pipenv && pipenv install django

# Expor a porta 8000 para acessar a aplicação
EXPOSE 8000

# Rodar a aplicação Django no servidor de desenvolvimento
CMD ["tail", "-f", "/dev/null"]
