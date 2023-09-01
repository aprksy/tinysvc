DATE=$(date '+%Y%m%d')
if [ $1 != "" ]; then
    DATE=$1
fi

source ./path.sh $DATE
DMIN1=$(date --date "$DATE -1 days" '+%Y%m%d')
DMIN2=$(date --date "$DATE -2 days" '+%Y%m%d')

mkdir -p $RESULTS

echo "STEP 1: preparing IRIIS data"
$SCRIPT_PATH/01_prepare_iriis.sh $DATE

echo "STEP 2: preparing IRONMAN data"
$SCRIPT_PATH/02_prepare_ironman.sh $DATE

echo "STEP 3: collect IRIIS cell based on IRONMAN data"
$SCRIPT_PATH/03_get_iriis_cells.sh $DATE

echo "STEP 4: collect Mediation GH7 cells based on IRIIS data"
$SCRIPT_PATH/04_get_mediation_cells.sh $DATE $DMIN1 0
$SCRIPT_PATH/04_get_mediation_cells.sh $DATE $DMIN2 0

echo "STEP 5: join IRIIS and Mediation cells"
$SCRIPT_PATH/05_join_iriis_mediation.sh $DATE $DMIN1 0
$SCRIPT_PATH/05_join_iriis_mediation.sh $DATE $DMIN2 0

echo "STEP 6: calculate ratio based on IRIIS definition"
$SCRIPT_PATH/06_calculate.sh $DATE $DMIN1 0 1
$SCRIPT_PATH/06_calculate.sh $DATE $DMIN2 0 2

echo "STEP 7: combine result between day-min-1 & day-min-2"
$SCRIPT_PATH/07_combine_days.sh $DATE 0 1,2

echo "STEP 8: combine result from all regionals"
$SCRIPT_PATH/08_combine_regions.sh $DATE

echo "STEP 9: delete old workspace data"
rm -rf $SCRIPT_PATH/$(date --date "$DATE -$WORKSPACE_RETENTION days" '+%Y%m%d')