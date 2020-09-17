docker stop $(docker ps -a -q) >/dev/null 2>&1
docker rm $(docker ps -a -q) >/dev/null 2>&1