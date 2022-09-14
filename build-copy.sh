
# build image
docker build . -t vinyl-player
# copy over to raspberry
docker save vinyl-player:latest | bzip2 | ssh titan docker load
