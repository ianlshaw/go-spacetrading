#!/usr/bin/bash
ls TVRJ.token
if [ $? -ne 0 ]
then
  echo "Local state does not exist"
  exit 0
else
  aws s3 cp TVRJ.markets.json s3://go-spacetraders-state/TVRJ.markets.json 
  aws s3 cp TVRJ.shipyards.json s3://go-spacetraders-state/TVRJ.shipyards.json 
  # careful with this one
  aws s3 cp TVRJ.token s3://go-spacetraders-state/TVRJ.token
  aws s3 cp TVRJ.waypoints.json s3://go-spacetraders-state/TVRJ.waypoints.json
  aws s3 cp TVRJ.world.json s3://go-spacetraders-state/TVRJ.world.json
fi