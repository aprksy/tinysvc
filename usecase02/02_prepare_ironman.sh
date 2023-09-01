source ./path.sh $1
DATE=$1

mkdir -p $IRONMAN_CELLS
tail +2 $IRONMAN_FILENAME \
    | tr -s '|' , \
    | cut -d, -f1,2,3 \
    | sort -t, -k2 -u \
    | sed 's/,/-/2' \
    > $IRONMAN_CELLS/cells.csv
echo "  cell count: "$(cat $IRONMAN_CELLS/cells.csv | wc -l)