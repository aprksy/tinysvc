DATE=$(date '+%Y%m%d')
if [ $1 != "" ]; then
    DATE=$1
fi

source ./path.sh $DATE
DMIN1=$(date --date "$DATE -1 days" '+%Y%m%d')
DMIN2=$(date --date "$DATE -2 days" '+%Y%m%d')

mkdir -p $RESULTS

echo "STEP 1: preparing CROSSFEEDER data"
echo "  preparing data on date $DMIN1"
$SCRIPT_PATH/01_prepare_crossfeeder.sh $DATE $DMIN1
echo "  preparing data on date $DMIN2"
$SCRIPT_PATH/01_prepare_crossfeeder.sh $DATE $DMIN2

echo "STEP 2: preparing IRONMAN data"
$SCRIPT_PATH/02_prepare_ironman.sh $DATE

echo "STEP 3: join IRONMAN vs CROSSFEEDER data"
$SCRIPT_PATH/03_join_ironman_crossfeeder.sh $DATE $DMIN1 $DMIN2

echo "STEP 4: delete old workspace data"
rm -rf $SCRIPT_PATH/$(date --date "$DATE -$WORKSPACE_RETENTION days" '+%Y%m%d')