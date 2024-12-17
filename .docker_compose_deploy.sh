zip -r dickobot.zip .
scp "dickobot.zip" ximanager@xi:~/
ssh ximanager@xi
unzip dickobot.zip -d dickobot
cd dickobot/
sudo docker compose build
sudo docker compose down
sudo docker compose up -dzip -r dickobot.zip . -x 'data/' 'dickobot.zip' 'dickobot/'