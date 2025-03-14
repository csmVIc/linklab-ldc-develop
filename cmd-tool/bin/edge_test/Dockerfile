FROM 47.111.77.0:32333/arm32v7/python:3.7-slim

WORKDIR /app

COPY . .

RUN apt-get update
RUN apt-get install -y wget
RUN apt-get install -y libopenblas-dev
RUN apt-get install -y gcc
RUN apt-get install -y libjpeg-dev
RUN apt-get install -y zlib1g-dev

RUN wget https://github.com/KumaTea/pytorch-arm/releases/download/v1.6.0/torch-1.6.0-cp37-cp37m-linux_armv7l.whl
RUN wget https://github.com/KumaTea/pytorch-arm/releases/download/v1.6.0/torchvision-0.7.0-cp37-cp37m-linux_armv7l.whl

RUN pip install torch-1.6.0-cp37-cp37m-linux_armv7l.whl
RUN pip install torchvision-0.7.0-cp37-cp37m-linux_armv7l.whl
RUN pip install --no-cache-dir -r requirements.txt
