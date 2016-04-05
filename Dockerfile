FROM golang:latest 
RUN apt-get update
RUN rm -rf /var/lib/apt/lists/*
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o hello . 
CMD ["/app/hello"]
