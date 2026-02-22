#!/usr/bin/bash

aws s3 ls s3://go-spacetraders-state/TVRJ.token
if [ $? -ne 0 ]
then
  echo "Remote state does not exist"
  exit 0
else
  # careful with this one
  aws s3 cp s3://go-spacetraders-state/TVRJ.token TVRJ.token 
  aws s3 cp s3://go-spacetraders-state/TVRJ.world.json TVRJ.world.json
fi