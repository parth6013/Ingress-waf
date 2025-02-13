FROM containers.cisco.com/ngcs/hcps-python:develop

WORKDIR /app

COPY requirements.txt /app/requirements.txt

RUN pip install -r requirements.txt --no-cache-dir && \
    rm /app/requirements.txt

COPY static  /app/static
COPY templates /app/templates
COPY app.py /app
COPY db.py /app
COPY database.db /app
COPY shell.sh /app

CMD ["python", "app.py"]
