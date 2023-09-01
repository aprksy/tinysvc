source ./path.sh $1
DATE=$1
DATE1=$2
DATE2=$3

mkdir -p $JOIN_IRONMAN_CROSSFEEDER_DIR
echo "  processing data $DATE1"
join -t, -1 2 -2 1 -a1 -e2 -o 1.1,0,2.2 \
    $IRONMAN_CELLS/cells.csv \
    $CROSSFEEDER_CELLS/cells-$DATE1.csv \
    > $JOIN_IRONMAN_CROSSFEEDER_DIR/cells1.csv

echo "  processing data $DATE2"
join -t, -1 2 -2 1 -a1 -e2 -o 1.1,0,1.3,2.2 \
    $JOIN_IRONMAN_CROSSFEEDER_DIR/cells1.csv \
    $CROSSFEEDER_CELLS/cells-$DATE2.csv \
    > $JOIN_IRONMAN_CROSSFEEDER_DIR/cells2.csv

cat $RESULT_COLS \
    | tr -s \\n , \
    | sed -E 's/,$/\n/1' \
    > $JOIN_IRONMAN_CROSSFEEDER_DIR/$DATE-crossfeeder.csv

echo "  merging data for on both dates"
cat $JOIN_IRONMAN_CROSSFEEDER_DIR/cells2.csv \
    | awk -F, 'BEGIN{OFS=","}{
        f1=$3; f2=$4;
        f3=$3 < $4 ? $3 : $4;
        s1=f1==0?"PASS":f1==1?"NOTPASS":"UNDEFINED";
        s2=f2==0?"PASS":f2==1?"NOTPASS":"UNDEFINED";
        s3=f3==0?"PASS":f3==1?"NOTPASS":"UNDEFINED";
        print $1,$2,s1,s2,s3
    }' \
    >> $JOIN_IRONMAN_CROSSFEEDER_DIR/$DATE-crossfeeder.csv

cp $JOIN_IRONMAN_CROSSFEEDER_DIR/$DATE-crossfeeder.csv $RESULTS/
    
rm $JOIN_IRONMAN_CROSSFEEDER_DIR/{cells1.csv,cells2.csv}