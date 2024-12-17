sudo docker build -t dickobot .
sudo docker kill dickobot
sudo docker run --detach --name dickobot -v $(pwd)/data/bots/dickobot:/data 'dickobot'