#!/usr/bin/bash

aws s3 ls s3://go-spacetraders-state/TVRJ.token
if [ $? -ne 0 ]
then
  echo "Remote state does not exist"
  exit 0
else
  aws s3 cp s3://go-spacetraders-state/TVRJ.markets.json TVRJ.markets.json 
  aws s3 cp s3://go-spacetraders-state/TVRJ.shipyards.json TVRJ.shipyards.json 
  # careful with this one
  aws s3 cp s3://go-spacetraders-state/TVRJ.token TVRJ.token 
  aws s3 cp s3://go-spacetraders-state/TVRJ.trade_routes TVRJ.trade_routes 
  aws s3 cp s3://go-spacetraders-state/TVRJ.waypoints.json TVRJ.waypoints.json
  aws s3 cp s3://go-spacetraders-state/TVRJ.world.json TVRJ.world.json
fi