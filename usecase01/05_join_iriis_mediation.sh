source ./path.sh $1
DATE=$1
DATE1=$2
REGIONAL=$3

mkdir -p $JOIN_IRIIS_GEOHASH_DIR
echo "  processing data $DATE1"
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    if [ $REGIONAL = $reg ] || [ $REGIONAL = "0" ]; then
        echo "  - joining cells in regional $reg"
        join -t, -j1 -e0 -a1 -a2 -o 0,1.2,2.2 \
            $IRIIS_CELLS/tiles_R$reg-defined.csv \
            $GEOHASH_DIR/tiles_${DATE1}_R$reg-defined.csv \
            > $JOIN_IRIIS_GEOHASH_DIR/temp.csv

        sed 's/-/,/2' $JOIN_IRIIS_GEOHASH_DIR/temp.csv \
            > $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R$reg-defined.csv
        rm $JOIN_IRIIS_GEOHASH_DIR/temp.csv

        rowcount=$(cat $GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv | wc -l)
        paste -d, \
            <(cut -d, -f1 $GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv) \
            <(yes 0 | head -$rowcount) \
            <(cut -d, -f2 $GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv) \
            | sed 's/-/,/2' \
            > $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv

        cat \
            $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R$reg-defined.csv \
            $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R$reg-undefined.csv \
            > $JOIN_IRIIS_GEOHASH_DIR/temp.csv
        mv \
            $JOIN_IRIIS_GEOHASH_DIR/temp.csv \
            $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R$reg.csv 
    fi
done