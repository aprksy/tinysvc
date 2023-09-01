source ./path.sh $1
DATE=$1
DATE1=$2
REGIONAL=$3

mkdir -p $GEOHASH_DIR
echo "  processing data $DATE1"
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    if [ $REGIONAL = $reg ] || [ $REGIONAL = "0" ]; then
        echo "    processing regional $reg"
        filename=$GEOHASH_RAW_DIR/$(echo $GEOHASH_RAW_FILENAME \
            | sed "s/REGION/$(seq -f '%02g' $reg $reg)/" \
            | sed "s/DATE/$DATE1/")
        paste -d, \
            <(cut -d, -f3,4 $filename 2> /dev/null) \
            <(cut -d, -f2,6 $filename 2> /dev/null) \
            | sed 's/,/-/1' \
            > $GEOHASH_DIR/tiles_${DATE1}_R$reg.csv
        echo "    - total tiles: "$(cat $GEOHASH_DIR/tiles_${DATE1}_R$reg.csv 2>/dev/null | wc -l)

        # process iriis-defined data
        grep -F -f $IRIIS_CELLS/cells_R$reg-defined.csv $GEOHASH_DIR/tiles_${DATE1}_R$reg.csv \
            | sort -t, -k1 \
            > $GEOHASH_DIR/temp.csv
        sed 's/,/-/1' $GEOHASH_DIR/temp.csv \
            | grep -v '1--1' \
            | grep -v '0-0' \
            > $GEOHASH_DIR/tiles_${DATE1}_R$reg-defined.csv
        rm $GEOHASH_DIR/temp.csv
        echo "    - tiles from defined cells: "$(cat $GEOHASH_DIR/tiles_${DATE1}_R$reg-defined.csv 2>/dev/null | wc -l)

        # process iriis-undefined data
        grep -F -f $IRIIS_CELLS/cells_R$reg-undefined.csv $GEOHASH_DIR/tiles_${DATE1}_R$reg.csv \
            | sort -t, -k1 \
            > $GEOHASH_DIR/temp.csv
        sed 's/,/-/1' $GEOHASH_DIR/temp.csv \
            > $GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv
        rm $GEOHASH_DIR/temp.csv
        echo "    - tiles from undefined cells: "$(cat $GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv 2>/dev/null | wc -l)
    fi
done