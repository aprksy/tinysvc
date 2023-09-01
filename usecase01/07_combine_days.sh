source ./path.sh $1
DATE=$1
REGIONAL=$2
DAYS=$3

mkdir -p $COMBINED_DAYS_DIR
echo "  processing data $DATE1"
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    if [ $REGIONAL = $reg ] || [ $REGIONAL = "0" ]; then
        echo "  - combining day-min-1 and day-min-2 for regional $reg"
        DAYS=$(echo $DAYS | tr -s , \\n)

        paste -d, \
                <(yes $reg \
                    | head \
                        -$(cat $CALCULATION_RESULT_DIR/result_R${reg}_D1.csv \
                        | wc -l)) \
                <(cut -d, -f1,2,3,4,5,6,7,8 \
                    $CALCULATION_RESULT_DIR/result_R${reg}_D1.csv) \
                <(cut -d, -f4,5,6,7,8 \
                    $CALCULATION_RESULT_DIR/result_R${reg}_D2.csv) \
                > $COMBINED_DAYS_DIR/result_R$reg.csv

        paste -d, \
                <(yes $reg \
                    | head \
                        -$(cat $CALCULATION_RESULT_DIR/result_R${reg}_D1-badgrids.csv \
                        | wc -l)) \
                <(cut -d, -f1,2,3,4 \
                    $CALCULATION_RESULT_DIR/result_R${reg}_D1-badgrids.csv) \
                <(cut -d, -f4 \
                    $CALCULATION_RESULT_DIR/result_R${reg}_D2-badgrids.csv) \
                > $COMBINED_DAYS_DIR/result_R$reg-badgrids.csv
    fi
done