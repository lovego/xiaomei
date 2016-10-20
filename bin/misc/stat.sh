#!/bin/bash
date=2016-07-21

bin/tools production run "egrep -ah '$date .* /plan/index/?[ ?#]' release/log/app.log" | fgrep /plan/index > tmp/statistics/plan_$date.txt

bin/tools production run "egrep -ah '$date .* /api/category_plan/cxdata/' release/log/app.log" | fgrep /api/category_plan/ > tmp/statistics/plan_app_$date.txt

bin/tools production run "egrep -ah '$date .* /vendor/plan/?[ ?#]' release/log/app.log" | fgrep /vendor/plan > tmp/statistics/vendor_plan_$date.txt

bin/tools production run "egrep -ah '$date .* /kpi/?[ ?#]' release/log/app.log" | fgrep /kpi > tmp/statistics/kpi_$date.txt

fgrep -a $date release/log/app.log | grep -Po '(/[^/? ]*[a-z][^/? ]*)+' | sort | uniq -c | sort -nrk1
