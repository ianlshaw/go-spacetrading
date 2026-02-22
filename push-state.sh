#!/usr/bin/bash
ls TVRJ.token
if [ $? -ne 0 ]
then
  echo "Local state does not exist"
  exit 0
else
  # careful with this one
  aws s3 cp TVRJ.token s3://go-spacetraders-state/TVRJ.token
  aws s3 cp TVRJ.world.json s3://go-spacetraders-state/TVRJ.world.json
fi