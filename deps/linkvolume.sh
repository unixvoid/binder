#!/bin/bash
DATA_DIR=binder-stack/binderdata

# check if symlink exists
if ! [ -L "$DATA_DIR" ]; then
  # create volume
  sudo docker volume create --name binderstack_datavolume
  VOLUME_PATH=$(sudo docker volume inspect --format '{{ .Mountpoint }}' binderstack_datavolume)
  
  # link volume
  sudo ln -s $VOLUME_PATH $DATA_DIR
fi
